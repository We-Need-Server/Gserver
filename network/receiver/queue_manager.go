package receiver

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/packet"
)

type QueueManager struct {
	queue *internal_type.Queue[packet.PacketI]
}

func NewQueueManager(queue *internal_type.Queue[packet.PacketI]) *QueueManager {
	return &QueueManager{queue: queue}
}

func (qm *QueueManager) EnManagedQueue(p packet.PacketI) {
	qm.queue.Enqueue(p)
}
