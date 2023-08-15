package transaction

import (
	"golang.org/x/xerrors"
)

func (m *Manager) initTronAddress(as []string) {
	err := m.SaveRechargeAddress(as)
	if err != nil {
		log.Errorf("SaveRechargeAddress err:%s", err.Error())
	}

	list, err := m.GetAllRechargeAddresses()
	if err != nil {
		log.Errorf("GetAllRechargeAddresses err:%s", err.Error())
	}

	for _, addr := range list {
		m.tronAddrs[addr.Addr] = addr.UserAddr
	}
}

func (m *Manager) addTronAddr(addr, userID string) {
	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	m.tronAddrs[addr] = userID
}

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
