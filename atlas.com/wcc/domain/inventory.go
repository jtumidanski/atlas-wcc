package domain

type Inventory struct {
   equipInventory EquipInventory
   useInventory   ItemInventory
   setupInventory ItemInventory
   etcInventory   ItemInventory
   cashInventory  ItemInventory
}

func EmptyInventory() Inventory {
   return Inventory{
      equipInventory: EmptyEquipInventory(),
      useInventory:   EmptyItemInventory(),
      setupInventory: EmptyItemInventory(),
      etcInventory:   EmptyItemInventory(),
      cashInventory:  EmptyItemInventory(),
   }
}

func (i Inventory) SetEquipInventory(ei EquipInventory) Inventory {
   return Inventory{
      equipInventory: ei,
      useInventory:   i.useInventory,
      setupInventory: i.setupInventory,
      etcInventory:   i.etcInventory,
      cashInventory:  i.cashInventory,
   }
}

func (i Inventory) SetUseInventory(ii ItemInventory) Inventory {
   return Inventory{
      equipInventory: i.equipInventory,
      useInventory:   ii,
      setupInventory: i.setupInventory,
      etcInventory:   i.etcInventory,
      cashInventory:  i.cashInventory,
   }
}

func (i Inventory) SetSetupInventory(ii ItemInventory) Inventory {
   return Inventory{
      equipInventory: i.equipInventory,
      useInventory:   i.useInventory,
      setupInventory: ii,
      etcInventory:   i.etcInventory,
      cashInventory:  i.cashInventory,
   }
}

func (i Inventory) SetEtcInventory(ii ItemInventory) Inventory {
   return Inventory{
      equipInventory: i.equipInventory,
      useInventory:   i.useInventory,
      setupInventory: i.setupInventory,
      etcInventory:   ii,
      cashInventory:  i.cashInventory,
   }
}

func (i Inventory) SetCashInventory(ii ItemInventory) Inventory {
   return Inventory{
      equipInventory: i.equipInventory,
      useInventory:   i.useInventory,
      setupInventory: i.setupInventory,
      etcInventory:   i.etcInventory,
      cashInventory:  ii,
   }
}

func (i Inventory) EquipInventory() EquipInventory {
   return i.equipInventory
}

func (i Inventory) UseInventory() ItemInventory {
   return i.useInventory
}

func (i Inventory) SetupInventory() ItemInventory {
   return i.setupInventory
}

func (i Inventory) EtcInventory() ItemInventory {
   return i.etcInventory
}

func (i Inventory) CashInventory() ItemInventory {
   return i.cashInventory
}
