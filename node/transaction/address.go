package transaction

func (m *Manager) initTronAddress(as []string) {
	err := m.SaveRechargeAddress(as)
	if err != nil {
		log.Errorf("SaveRechargeAddress err:%s", err.Error())
	}

	list, err := m.LoadUsedRechargeAddresses()
	if err != nil {
		log.Errorf("GetAllRechargeAddresses err:%s", err.Error())
	}

	for _, addr := range list {
		m.tronAddrs[addr.Addr] = addr.UserID
	}
}

func (m *Manager) addTronAddr(addr, userID string) {
	m.tronAddrLock.Lock()
	defer m.tronAddrLock.Unlock()

	m.tronAddrs[addr] = userID
}
