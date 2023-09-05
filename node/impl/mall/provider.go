package mall

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/LMF709268224/titan-vps/node/utils"
)

func (m *Mall) SetAccessKeyInfo(ctx context.Context, info *types.AccessKeyInfo) error {
	providerID := handler.GetID(ctx)

	info.ProviderID = providerID

	return m.SaveAccessKeyInfo(info)
}

func (m *Mall) UpdateAccessKeyNick(ctx context.Context, accessSecret, nickName string) error {
	providerID := handler.GetID(ctx)

	info, err := m.LoadAccessKeyInfo(providerID, accessSecret)
	if err != nil {
		return err
	}

	info.NickName = nickName

	return m.UpdateAccessKeyInfo(info)
}

func (m *Mall) RemoveAccessKeyInfo(ctx context.Context, accessSecret string) error {
	providerID := handler.GetID(ctx)

	return m.DeleteAccessKeyInfo(providerID, accessSecret)
}

// GetInstanceRecordsOfAccessKey
func (m *Mall) GetInstanceRecordsOfAccessKey(ctx context.Context, limit, page int64, accessSecret string) (*types.GetInstanceResponse, error) {
	out := &types.GetInstanceResponse{}

	rows, total, err := m.LoadInstancesInfoByAccessKey(accessSecret, limit, page)
	if err != nil {
		log.Errorf("LoadInstancesInfoByAccessKey err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}
	defer rows.Close()

	out.Total = total

	for rows.Next() {
		info := &types.InstanceDetails{}
		err = rows.StructScan(info)
		if err != nil {
			log.Errorf("InstanceDetails StructScan err: %s", err.Error())
			continue
		}

		orders, err := m.LoadOrderRecordsByVpsID(info.ID, types.Done, types.OrderDoneStateSuccess)
		if err == nil {
			tradePrice := "0"
			for _, order := range orders {
				price, err := utils.AddBigInt(tradePrice, order.Value)
				if err != nil {
					log.Errorf("AddBigInt %s,%s ,%s", tradePrice, order.Value, err.Error())
					continue
				}

				tradePrice = price
			}

			info.Value = tradePrice
		}

		out.List = append(out.List, m.VpsMgr.UpdateInstanceInfo(info, false))
	}

	return out, nil
}
