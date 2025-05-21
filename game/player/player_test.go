package player

//import (
//	"testing"
//)
//
//func TestNewPlayer(t *testing.T) {
//	p := NewPlayer()
//	if p == nil {
//		t.Error("NewPlayer() returned nil")
//	}
//
//	// 기본값 테스트 (모든 필드가 0으로 초기화되어야 함)
//	if p.PositionX != 0 || p.PositionY != 0 || p.PositionZ != 0 ||
//		p.XDelta != 0 || p.YDelta != 0 || p.ZDelta != 0 ||
//		p.YawAngle != 0 || p.YawDelta != 0 ||
//		p.PTAngle != 0 || p.PTDelta != 0 {
//		t.Error("NewPlayer() did not initialize with zero values")
//	}
//}
//
//func TestGetPlayerInfo(t *testing.T) {
//	// 테스트용 Player 객체 설정
//	original := Player{
//		PositionX: 10.5,
//		XDelta:    1.5,
//		PositionY: 20.5,
//		YDelta:    2.5,
//		PositionZ: 30.5,
//		ZDelta:    3.5,
//		YawAngle:  45.0,
//		YawDelta:  5.0,
//		PTAngle:   90.0,
//		PTDelta:   10.0,
//	}
//
//	// GetPlayerInfo 메서드 호출
//	info := original.GetPlayerInfo()
//
//	// 반환된 정보가 원본과 동일한지 확인
//	if info.PositionX != original.PositionX ||
//		info.XDelta != original.XDelta ||
//		info.PositionY != original.PositionY ||
//		info.YDelta != original.YDelta ||
//		info.PositionZ != original.PositionZ ||
//		info.ZDelta != original.ZDelta ||
//		info.YawAngle != original.YawAngle ||
//		info.YawDelta != original.YawDelta ||
//		info.PTAngle != original.PTAngle ||
//		info.PTDelta != original.PTDelta {
//		t.Error("GetPlayerInfo() did not return the correct player information")
//	}
//}
//
//func TestMoveFoward(t *testing.T) {
//	p := NewPlayer()
//
//	// 테스트 케이스 1: 양수 delta 값
//	p.MoveForward(5.0)
//	if p.XDelta != 5.0 {
//		t.Errorf("MoveFoward(5.0) failed, expected XDelta to be 5.0, got %v", p.XDelta)
//	}
//
//	// 테스트 케이스 2: 음수 delta 값
//	p.MoveForward(-2.0)
//	if p.XDelta != 3.0 { // 5.0 + (-2.0) = 3.0
//		t.Errorf("MoveFoward(-2.0) failed, expected XDelta to be 3.0, got %v", p.XDelta)
//	}
//
//	// 테스트 케이스 3: 0 delta 값
//	p.MoveForward(0.0)
//	if p.XDelta != 3.0 { // 값이 변하지 않아야 함
//		t.Errorf("MoveFoward(0.0) failed, expected XDelta to remain 3.0, got %v", p.XDelta)
//	}
//}
//
//func TestMoveSide(t *testing.T) {
//	p := NewPlayer()
//
//	// 테스트 케이스 1: 양수 delta 값
//	p.MoveSide(8.0)
//	if p.ZDelta != 8.0 {
//		t.Errorf("MoveSide(8.0) failed, expected ZDelta to be 8.0, got %v", p.ZDelta)
//	}
//
//	// 테스트 케이스 2: 음수 delta 값
//	p.MoveSide(-3.0)
//	if p.ZDelta != 5.0 { // 8.0 + (-3.0) = 5.0
//		t.Errorf("MoveSide(-3.0) failed, expected ZDelta to be 5.0, got %v", p.ZDelta)
//	}
//
//	// 테스트 케이스 3: 0 delta 값
//	p.MoveSide(0.0)
//	if p.ZDelta != 5.0 { // 값이 변하지 않아야 함
//		t.Errorf("MoveSide(0.0) failed, expected ZDelta to remain 5.0, got %v", p.ZDelta)
//	}
//}
//
//func TestTransferYaw(t *testing.T) {
//	p := NewPlayer()
//
//	// 테스트 케이스 1: 양수 delta 값
//	p.TransferYaw(45.0)
//	if p.YawDelta != 45.0 {
//		t.Errorf("TransferYaw(45.0) failed, expected YawDelta to be 45.0, got %v", p.YawDelta)
//	}
//
//	// 테스트 케이스 2: 음수 delta 값
//	p.TransferYaw(-15.0)
//	if p.YawDelta != 30.0 { // 45.0 + (-15.0) = 30.0
//		t.Errorf("TransferYaw(-15.0) failed, expected YawDelta to be 30.0, got %v", p.YawDelta)
//	}
//
//	// 테스트 케이스 3: 0 delta 값
//	p.TransferYaw(0.0)
//	if p.YawDelta != 30.0 { // 값이 변하지 않아야 함
//		t.Errorf("TransferYaw(0.0) failed, expected YawDelta to remain 30.0, got %v", p.YawDelta)
//	}
//}
//
//func TestTransferPT(t *testing.T) {
//	p := NewPlayer()
//
//	// 테스트 케이스 1: 양수 delta 값
//	p.TransferPT(90.0)
//	if p.PTDelta != 90.0 {
//		t.Errorf("TransferPT(90.0) failed, expected PTDelta to be 90.0, got %v", p.PTDelta)
//	}
//
//	// 테스트 케이스 2: 음수 delta 값
//	p.TransferPT(-30.0)
//	if p.PTDelta != 60.0 { // 90.0 + (-30.0) = 60.0
//		t.Errorf("TransferPT(-30.0) failed, expected PTDelta to be 60.0, got %v", p.PTDelta)
//	}
//
//	// 테스트 케이스 3: 0 delta 값
//	p.TransferPT(0.0)
//	if p.PTDelta != 60.0 { // 값이 변하지 않아야 함
//		t.Errorf("TransferPT(0.0) failed, expected PTDelta to remain 60.0, got %v", p.PTDelta)
//	}
//}
//
//// 통합 테스트: 여러 메서드를 연속적으로 호출
//func TestIntegration(t *testing.T) {
//	p := NewPlayer()
//
//	// 여러 동작 수행
//	p.MoveForward(10.0)
//	p.MoveSide(5.0)
//	p.TransferYaw(30.0)
//	p.TransferPT(60.0)
//
//	// 결과 검증
//	if p.XDelta != 10.0 || p.ZDelta != 5.0 || p.YawDelta != 30.0 || p.PTDelta != 60.0 {
//		t.Errorf("Integration test failed, got XDelta=%v, ZDelta=%v, YawDelta=%v, PTDelta=%v",
//			p.XDelta, p.ZDelta, p.YawDelta, p.PTDelta)
//	}
//
//	// 추가 동작 수행 (누적 검증)
//	p.MoveForward(-5.0)
//	p.MoveSide(2.0)
//
//	// 결과 검증
//	if p.XDelta != 5.0 || p.ZDelta != 7.0 {
//		t.Errorf("Integration test (additional operations) failed, got XDelta=%v, ZDelta=%v",
//			p.XDelta, p.ZDelta)
//	}
//}
