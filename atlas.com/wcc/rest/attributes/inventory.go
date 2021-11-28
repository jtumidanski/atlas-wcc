package attributes

import (
	"atlas-wcc/rest/response"
	"encoding/json"
	"strconv"
)

const (
	ItemAttributesType      string = "com.atlas.cos.rest.attribute.ItemAttributes"
	EquipmentAttributesType string = "com.atlas.cos.rest.attribute.EquipmentAttributes"
	EquipmentStatisticsType string = "com.atlas.cos.rest.attribute.EquipmentStatisticsAttributes"
)

var equipmentIncludes = []response.ConditionalMapperProvider{
	transformItemAttributes,
	transformEquipmentAttributes,
	transformEquipmentStatistics,
}

type InventoryDataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type InventoryData struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes InventoryAttributes `json:"attributes"`
}

type InventoryAttributes struct {
	Type     string `json:"type"`
	Capacity byte   `json:"capacity"`
}

func (c *InventoryDataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	if len(c.included) == 1 {
		t.Included = c.included[0]
	} else {
		t.Included = c.included
	}
	return json.Marshal(t)
}

func (c *InventoryDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyInventoryData), equipmentIncludes...)
	if err != nil {
		return err
	}

	c.data = d
	c.included = i
	return nil
}

func (c *InventoryDataContainer) Data() *InventoryData {
	if len(c.data) >= 1 {
		return c.data[0].(*InventoryData)
	}
	return nil
}

func (c *InventoryDataContainer) DataList() []InventoryData {
	var r = make([]InventoryData, 0)
	for _, x := range c.data {
		r = append(r, *x.(*InventoryData))
	}
	return r
}

func (c *InventoryDataContainer) GetIncludedEquippedItems() []EquipmentData {
	var e = make([]EquipmentData, 0)
	for _, x := range c.included {
		if val, ok := x.(*EquipmentData); ok && val.Attributes.Slot < 0 {
			e = append(e, *val)
		}
	}
	return e
}

func (c *InventoryDataContainer) GetIncludedEquips() []EquipmentData {
	var e = make([]EquipmentData, 0)
	for _, x := range c.included {
		if val, ok := x.(*EquipmentData); ok && val.Attributes.Slot >= 0 {
			e = append(e, *val)
		}
	}
	return e
}

func (c *InventoryDataContainer) GetEquipmentStatistics(id int) *EquipmentStatisticsAttributes {
	for _, x := range c.included {
		if val, ok := x.(*EquipmentStatisticsData); ok {
			eid, err := strconv.Atoi(val.Id)
			if err == nil && eid == id {
				return &val.Attributes
			}
		}
	}
	return nil
}

func (c *InventoryDataContainer) GetIncludedItems() []ItemData {
	var e = make([]ItemData, 0)
	for _, x := range c.included {
		if val, ok := x.(*ItemData); ok && val.Attributes.Slot >= 0 {
			e = append(e, *val)
		}
	}
	return e
}

func EmptyItemData() interface{} {
	return &ItemData{}
}

type ItemData struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes ItemAttributes `json:"attributes"`
}

type ItemAttributes struct {
	ItemId   uint32 `json:"itemId"`
	Quantity uint16 `json:"quantity"`
	Slot     int16  `json:"slot"`
}

func transformItemAttributes() (string, response.ObjectMapper) {
	return response.UnmarshalData(ItemAttributesType, EmptyItemData)
}

func transformEquipmentAttributes() (string, response.ObjectMapper) {
	return response.UnmarshalData(EquipmentAttributesType, EmptyEquipmentData)
}

func transformEquipmentStatistics() (string, response.ObjectMapper) {
	return response.UnmarshalData(EquipmentStatisticsType, EmptyEquipmentStatisticsData)
}

func EmptyInventoryData() interface{} {
	return &InventoryData{}
}
