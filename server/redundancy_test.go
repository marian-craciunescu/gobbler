package server

import (
	"time"

	"github.com/smancke/guble/testutil"
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_Subscribe_on_random_node(t *testing.T) {
	testutil.SkipIfShort(t)
	a := assert.New(t)

	node1 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8080",
		NodeID:     1,
		NodePort:   10000,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node1)
	defer node1.cleanup(true)

	node2 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8081",
		NodeID:     2,
		NodePort:   10001,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node2)
	defer node2.cleanup(true)

	node1.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)
	node2.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)

	// subscribe on first node
	node1.Subscribe(fcmTopic, "1")

	// connect a client and send a message
	client1, err := node1.client("user1", 1000, true)
	a.NoError(err)

	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// only one message should be received but only on the first node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(1)
	node2.FCM.checkReceived(0)
}

func Test_Subscribe_working_After_Node_Restart(t *testing.T) {
	// defer testutil.EnableDebugForMethod()()
	testutil.SkipIfShort(t)
	a := assert.New(t)

	node1 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8082",
		NodeID:     1,
		NodePort:   10000,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node1)

	node2 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8083",
		NodeID:     2,
		NodePort:   10001,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node2)
	defer node2.cleanup(true)

	node1.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)
	node2.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)

	// subscribe on first node
	node1.Subscribe(fcmTopic, "1")

	// connect a clinet and send a message
	client1, err := node1.client("user1", 1000, true)
	a.NoError(err)
	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// one message should be received but only on the first node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(1)
	node2.FCM.checkReceived(0)

	// stop a node, cleanup without removing directories
	node1.cleanup(false)
	time.Sleep(time.Millisecond * 150)

	// restart the service
	restartedNode1 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen:  ":8082",
		StoragePath: node1.StoragePath,
		NodeID:      1,
		NodePort:    10000,
		Remotes:     "0.0.0.0:10000",
	})
	a.NotNil(restartedNode1)
	defer restartedNode1.cleanup(true)

	restartedNode1.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)

	// send a message to the former subscription.
	client1, err = restartedNode1.client("user1", 1000, true)
	a.NoError(err)
	time.Sleep(time.Second)

	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err, "Subscription should work even after node restart")

	// only one message should be received but only on the first node.
	// Every message should be delivered only once.
	restartedNode1.FCM.checkReceived(1)
	node2.FCM.checkReceived(0)
}

func Test_Independent_Receiving(t *testing.T) {
	testutil.SkipIfShort(t)
	a := assert.New(t)

	node1 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8084",
		NodeID:     1,
		NodePort:   10000,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node1)
	defer node1.cleanup(true)

	node2 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8085",
		NodeID:     2,
		NodePort:   10001,
		Remotes:    "0.0.0.0:10000",
	})
	a.NotNil(node2)
	defer node2.cleanup(true)

	node1.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)
	node2.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)

	// subscribe on first node
	node1.Subscribe(fcmTopic, "1")

	// connect a clinet and send a message
	client1, err := node1.client("user1", 1000, true)
	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// only one message should be received but only on the first node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(1)
	node2.FCM.checkReceived(0)

	// reset the counter
	node1.FCM.reset()

	// NOW connect to second node
	client2, err := node2.client("user2", 1000, true)
	a.NoError(err)
	err = client2.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// only one message should be received but only on the second node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(0)
	node2.FCM.checkReceived(1)
}

func Test_NoReceiving_After_Unsubscribe(t *testing.T) {
	testutil.SkipIfShort(t)
	a := assert.New(t)

	node1 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8086",
		NodeID:     1,
		NodePort:   10002,
		Remotes:    "0.0.0.0:10002",
	})
	a.NotNil(node1)
	defer node1.cleanup(true)

	node2 := newTestClusterNode(t, testClusterNodeConfig{
		HttpListen: "0.0.0.0:8087",
		NodeID:     2,
		NodePort:   10003,
		Remotes:    "0.0.0.0:10002",
	})
	a.NotNil(node2)
	defer node2.cleanup(true)

	node1.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)
	node2.FCM.setupRoundTripper(20*time.Millisecond, 10, testutil.SuccessFCMResponse)

	// subscribe on first node
	node1.Subscribe(fcmTopic, "1")
	time.Sleep(50 * time.Millisecond)

	// connect a client and send a message
	client1, err := node1.client("user1", 1000, true)
	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// only one message should be received but only on the first node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(1)
	node2.FCM.checkReceived(0)

	// Unsubscribe
	node2.Unsubscribe(fcmTopic, "1")
	time.Sleep(50 * time.Millisecond)

	// reset the counter
	node1.FCM.reset()

	// and send a message again. No one should receive it
	err = client1.Send(fcmTopic, "body", "{jsonHeader:1}")
	a.NoError(err)

	// only one message should be received but only on the second node.
	// Every message should be delivered only once.
	node1.FCM.checkReceived(0)
	node2.FCM.checkReceived(0)
}
