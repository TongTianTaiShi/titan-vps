package transaction

import (
	"golang.org/x/xerrors"
)

func (m *Manager) initFvmAddress(as []string) {
	defer m.addrWait.Done()

	m.fvmAddrLock.Lock()
	defer m.fvmAddrLock.Unlock()

	for _, addr := range as {
		m.usabilityFvmAddrs[addr] = ""
	}
}

func (m *Manager) initTronAddress(as []string) {
	defer m.addrWait.Done()

	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	for _, addr := range as {
		m.usabilityTronAddrs[addr] = ""
	}
}

// AllocateFvmAddress get a fvm address
func (m *Manager) AllocateFvmAddress(orderID string) (string, error) {
	m.fvmAddrLock.Lock()
	defer m.fvmAddrLock.Unlock()

	if len(m.usabilityFvmAddrs) > 0 {
		for addr := range m.usabilityFvmAddrs {
			m.usedFvmAddrs[addr] = orderID
			delete(m.usabilityFvmAddrs, addr)
			return addr, nil
		}
	}

	return "", xerrors.New("not found address")
}

// RevertFvmAddress revert a fvm address
func (m *Manager) RevertFvmAddress(addr string) {
	m.fvmAddrLock.Lock()
	defer m.fvmAddrLock.Unlock()

	delete(m.usedFvmAddrs, addr)
	m.usabilityFvmAddrs[addr] = ""
}

// RecoverOutstandingFvmOrders recover the fvm order
func (m *Manager) RecoverOutstandingFvmOrders(addr, orderID string) {
	m.addrWait.Wait()

	m.fvmAddrLock.Lock()
	defer m.fvmAddrLock.Unlock()

	m.usedFvmAddrs[addr] = orderID
	delete(m.usabilityFvmAddrs, addr)
}

// AllocateTronAddress get a tron address
func (m *Manager) AllocateTronAddress(orderID string) (string, error) {
	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	if len(m.usabilityTronAddrs) > 0 {
		for addr := range m.usabilityTronAddrs {
			m.usedTronAddrs[addr] = orderID
			delete(m.usabilityTronAddrs, addr)
			return addr, nil
		}
	}
	return "", xerrors.New("not found address")
}

// RevertTronAddress revert a tron address
func (m *Manager) RevertTronAddress(addr string) {
	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	delete(m.usedTronAddrs, addr)
	m.usabilityTronAddrs[addr] = ""
}

// RecoverOutstandingTronOrders recover the tron order
func (m *Manager) RecoverOutstandingTronOrders(addr, orderID string) {
	m.addrWait.Wait()

	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	m.usedTronAddrs[addr] = orderID
	delete(m.usabilityTronAddrs, addr)
}
