package receiver

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/packet"
)

type QueueManager struct {
	QChan chan *packet.PacketI
	queue *internal_type.Queue[*packet.PacketI]
}

func NewQueueManager(queue *internal_type.Queue[*packet.PacketI]) *QueueManager {
	return &QueueManager{QChan: make(chan *packet.PacketI), queue: queue}
}

func (qm *QueueManager) StartQueueManager() {
	for {
		select {
		case p := <-qm.QChan:
			// 채널에서 패킷을 받아서 큐에 추가
			if p != nil {
				qm.queue.Enqueue(p)
			}
			// 필요하다면 다른 케이스도 추가 가능
			// case <-ctx.Done():
			//    return // 종료 시그널이 필요한 경우
		}
	}
}

// origin EnQueue but, EnManagedQueue lol
func (qm *QueueManager) EnManagedQueue(p *packet.PacketI) {
	qm.queue.Enqueue(p)
}
