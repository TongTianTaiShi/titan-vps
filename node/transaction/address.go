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

// GetTronAddr get a fvm address
func (m *Manager) GetTronAddr() string {
	return m.tronAddr
}
