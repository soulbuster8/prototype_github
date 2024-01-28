package main

// This file will define the consistent hash ring where we will define the
// number of nodes which we wanted to add in the ring. Implementation of
// hash ring will be based upon the array.
// This file will also expose method to add node and delete node from the
// hash ring.
import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const modulo_number = 57

var node_id_map = make(map[int32]bool)

type ConsistentHashRing struct {
	hash_ring []int32

	total_nodes int32
}

//-------------------------------------------------------------------------------------------

func createConsistentHashingRing(total_nodes uint32) *ConsistentHashRing {
	var consistent_hash_ring ConsistentHashRing
	consistent_hash_ring.hash_ring = make([]int32, total_nodes)

	rand.Seed(time.Now().UnixNano())
	for ii := 0; ii < int(total_nodes); {
		node_id := rand.Int31n(modulo_number)
		_, exists := node_id_map[node_id]

		if exists {
			continue
		}

		node_id_map[node_id] = true
		consistent_hash_ring.hash_ring[ii] = node_id
		ii++
	}

	sort.Slice(consistent_hash_ring.hash_ring, func(idx1, idx2 int) bool {
		return consistent_hash_ring.hash_ring[idx1] < consistent_hash_ring.hash_ring[idx2]
	})

	return &consistent_hash_ring
}

//-------------------------------------------------------------------------------------------

// This function will add a node in the consistent hash ring.
func (consistent_hash_ring *ConsistentHashRing) addNode() {
	for true {
		node_id := rand.Int31n(modulo_number)
		_, exists := node_id_map[node_id]

		// If the key already exists in the hash ring. We should not
		// put it again.
		if exists {
			continue
		}

		fmt.Println("Adding node_id to the hash ring ", node_id)

		add_index := findGreaterThanEqualToKey(node_id, consistent_hash_ring.hash_ring)
		consistent_hash_ring.insertIntoHashRing(add_index, node_id)
		break
	}
}

//-------------------------------------------------------------------------------------------

// This function will give us the index where we have to insert the specified key.
// Assumption here is that slice which is passed as an argument is sorted and key
// is not present in the given slice as well.
func findGreaterThanEqualToKey(key int32, consistent_hash_ring []int32) int32 {
	var high_index int32 = int32(len(consistent_hash_ring) - 1)
	var low_index int32 = 0
	for low_index <= high_index {
		mid_index := low_index + (high_index-low_index)/2

		if consistent_hash_ring[mid_index] == key {
			return mid_index
		} else if consistent_hash_ring[mid_index] > key {
			high_index = mid_index - 1

			// It is possible that key which we want to insert is either on the 0th index or
			// it does lie between mid_index and mid_index - 1.
			if high_index < 0 || consistent_hash_ring[high_index] < key {
				return mid_index
			}
		} else if consistent_hash_ring[mid_index] < key {
			low_index = mid_index + 1

			// key might be inserted in the end of the array or it should be lying between mid_index
			// and mid_index + 1.
			if low_index >= int32(len(consistent_hash_ring)) || consistent_hash_ring[low_index] > key {
				return mid_index + 1
			}
		}
	}

	return low_index
}

//-------------------------------------------------------------------------------------------

func (consistent_hash_ring *ConsistentHashRing) insertIntoHashRing(index int32, value int32) {
	consistent_hash_ring.hash_ring = append(consistent_hash_ring.hash_ring, 0)

	copy(consistent_hash_ring.hash_ring[index+1:], consistent_hash_ring.hash_ring[index:])

	consistent_hash_ring.hash_ring[index] = value
}

//-------------------------------------------------------------------------------------------

func (consistent_hash_ring *ConsistentHashRing) deleteNode(node_id int32) {
	_, exists := node_id_map[node_id]

	if exists {
		panic("Node should be present in the map.")
	}

	deleted_index := findGreaterThanEqualToKey(node_id, consistent_hash_ring.hash_ring)

	consistent_hash_ring.deleteIntoHashRing(deleted_index)
}

//-------------------------------------------------------------------------------------------

func (consistent_hash_ring *ConsistentHashRing) deleteIntoHashRing(index int32) {
	consistent_hash_ring.hash_ring = append(consistent_hash_ring.hash_ring[:index],
		consistent_hash_ring.hash_ring[index+1:]...)
}

//-------------------------------------------------------------------------------------------

func (consistent_hash_ring *ConsistentHashRing) printConsistentHashRing() {
	fmt.Println("Printing consitent hash ring.")
	for ii := 0; ii < len(consistent_hash_ring.hash_ring); ii++ {
		fmt.Print(consistent_hash_ring.hash_ring[ii], " ")
	}
	fmt.Println("Consistent hash ring is printed.")
}

//-------------------------------------------------------------------------------------------

func main() {
	consistent_hash_ring := createConsistentHashingRing(5)

	consistent_hash_ring.printConsistentHashRing()

	consistent_hash_ring.addNode()
	consistent_hash_ring.printConsistentHashRing()

	consistent_hash_ring.deleteNode(1)
	consistent_hash_ring.printConsistentHashRing()

}
