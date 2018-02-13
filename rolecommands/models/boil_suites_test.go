// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("RoleCommands", testRoleCommands)
	t.Run("RoleGroups", testRoleGroups)
	t.Run("RoleMenuOptions", testRoleMenuOptions)
	t.Run("RoleMenus", testRoleMenus)
}

func TestDelete(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsDelete)
	t.Run("RoleGroups", testRoleGroupsDelete)
	t.Run("RoleMenuOptions", testRoleMenuOptionsDelete)
	t.Run("RoleMenus", testRoleMenusDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsQueryDeleteAll)
	t.Run("RoleGroups", testRoleGroupsQueryDeleteAll)
	t.Run("RoleMenuOptions", testRoleMenuOptionsQueryDeleteAll)
	t.Run("RoleMenus", testRoleMenusQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsSliceDeleteAll)
	t.Run("RoleGroups", testRoleGroupsSliceDeleteAll)
	t.Run("RoleMenuOptions", testRoleMenuOptionsSliceDeleteAll)
	t.Run("RoleMenus", testRoleMenusSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsExists)
	t.Run("RoleGroups", testRoleGroupsExists)
	t.Run("RoleMenuOptions", testRoleMenuOptionsExists)
	t.Run("RoleMenus", testRoleMenusExists)
}

func TestFind(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsFind)
	t.Run("RoleGroups", testRoleGroupsFind)
	t.Run("RoleMenuOptions", testRoleMenuOptionsFind)
	t.Run("RoleMenus", testRoleMenusFind)
}

func TestBind(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsBind)
	t.Run("RoleGroups", testRoleGroupsBind)
	t.Run("RoleMenuOptions", testRoleMenuOptionsBind)
	t.Run("RoleMenus", testRoleMenusBind)
}

func TestOne(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsOne)
	t.Run("RoleGroups", testRoleGroupsOne)
	t.Run("RoleMenuOptions", testRoleMenuOptionsOne)
	t.Run("RoleMenus", testRoleMenusOne)
}

func TestAll(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsAll)
	t.Run("RoleGroups", testRoleGroupsAll)
	t.Run("RoleMenuOptions", testRoleMenuOptionsAll)
	t.Run("RoleMenus", testRoleMenusAll)
}

func TestCount(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsCount)
	t.Run("RoleGroups", testRoleGroupsCount)
	t.Run("RoleMenuOptions", testRoleMenuOptionsCount)
	t.Run("RoleMenus", testRoleMenusCount)
}

func TestInsert(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsInsert)
	t.Run("RoleCommands", testRoleCommandsInsertWhitelist)
	t.Run("RoleGroups", testRoleGroupsInsert)
	t.Run("RoleGroups", testRoleGroupsInsertWhitelist)
	t.Run("RoleMenuOptions", testRoleMenuOptionsInsert)
	t.Run("RoleMenuOptions", testRoleMenuOptionsInsertWhitelist)
	t.Run("RoleMenus", testRoleMenusInsert)
	t.Run("RoleMenus", testRoleMenusInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("RoleCommandToRoleGroupUsingRoleGroup", testRoleCommandToOneRoleGroupUsingRoleGroup)
	t.Run("RoleMenuOptionToRoleCommandUsingRoleCommand", testRoleMenuOptionToOneRoleCommandUsingRoleCommand)
	t.Run("RoleMenuOptionToRoleMenuUsingRoleMenu", testRoleMenuOptionToOneRoleMenuUsingRoleMenu)
	t.Run("RoleMenuToRoleCommandUsingNextRoleCommand", testRoleMenuToOneRoleCommandUsingNextRoleCommand)
	t.Run("RoleMenuToRoleGroupUsingRoleGroup", testRoleMenuToOneRoleGroupUsingRoleGroup)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("RoleCommandToRoleMenuOptions", testRoleCommandToManyRoleMenuOptions)
	t.Run("RoleCommandToNextRoleCommandRoleMenus", testRoleCommandToManyNextRoleCommandRoleMenus)
	t.Run("RoleGroupToRoleCommands", testRoleGroupToManyRoleCommands)
	t.Run("RoleGroupToRoleMenus", testRoleGroupToManyRoleMenus)
	t.Run("RoleMenuToRoleMenuOptions", testRoleMenuToManyRoleMenuOptions)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("RoleCommandToRoleGroupUsingRoleGroup", testRoleCommandToOneSetOpRoleGroupUsingRoleGroup)
	t.Run("RoleMenuOptionToRoleCommandUsingRoleCommand", testRoleMenuOptionToOneSetOpRoleCommandUsingRoleCommand)
	t.Run("RoleMenuOptionToRoleMenuUsingRoleMenu", testRoleMenuOptionToOneSetOpRoleMenuUsingRoleMenu)
	t.Run("RoleMenuToRoleCommandUsingNextRoleCommand", testRoleMenuToOneSetOpRoleCommandUsingNextRoleCommand)
	t.Run("RoleMenuToRoleGroupUsingRoleGroup", testRoleMenuToOneSetOpRoleGroupUsingRoleGroup)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("RoleCommandToRoleGroupUsingRoleGroup", testRoleCommandToOneRemoveOpRoleGroupUsingRoleGroup)
	t.Run("RoleMenuOptionToRoleCommandUsingRoleCommand", testRoleMenuOptionToOneRemoveOpRoleCommandUsingRoleCommand)
	t.Run("RoleMenuToRoleCommandUsingNextRoleCommand", testRoleMenuToOneRemoveOpRoleCommandUsingNextRoleCommand)
	t.Run("RoleMenuToRoleGroupUsingRoleGroup", testRoleMenuToOneRemoveOpRoleGroupUsingRoleGroup)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("RoleCommandToRoleMenuOptions", testRoleCommandToManyAddOpRoleMenuOptions)
	t.Run("RoleCommandToNextRoleCommandRoleMenus", testRoleCommandToManyAddOpNextRoleCommandRoleMenus)
	t.Run("RoleGroupToRoleCommands", testRoleGroupToManyAddOpRoleCommands)
	t.Run("RoleGroupToRoleMenus", testRoleGroupToManyAddOpRoleMenus)
	t.Run("RoleMenuToRoleMenuOptions", testRoleMenuToManyAddOpRoleMenuOptions)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("RoleCommandToRoleMenuOptions", testRoleCommandToManySetOpRoleMenuOptions)
	t.Run("RoleCommandToNextRoleCommandRoleMenus", testRoleCommandToManySetOpNextRoleCommandRoleMenus)
	t.Run("RoleGroupToRoleCommands", testRoleGroupToManySetOpRoleCommands)
	t.Run("RoleGroupToRoleMenus", testRoleGroupToManySetOpRoleMenus)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("RoleCommandToRoleMenuOptions", testRoleCommandToManyRemoveOpRoleMenuOptions)
	t.Run("RoleCommandToNextRoleCommandRoleMenus", testRoleCommandToManyRemoveOpNextRoleCommandRoleMenus)
	t.Run("RoleGroupToRoleCommands", testRoleGroupToManyRemoveOpRoleCommands)
	t.Run("RoleGroupToRoleMenus", testRoleGroupToManyRemoveOpRoleMenus)
}

func TestReload(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsReload)
	t.Run("RoleGroups", testRoleGroupsReload)
	t.Run("RoleMenuOptions", testRoleMenuOptionsReload)
	t.Run("RoleMenus", testRoleMenusReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsReloadAll)
	t.Run("RoleGroups", testRoleGroupsReloadAll)
	t.Run("RoleMenuOptions", testRoleMenuOptionsReloadAll)
	t.Run("RoleMenus", testRoleMenusReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsSelect)
	t.Run("RoleGroups", testRoleGroupsSelect)
	t.Run("RoleMenuOptions", testRoleMenuOptionsSelect)
	t.Run("RoleMenus", testRoleMenusSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsUpdate)
	t.Run("RoleGroups", testRoleGroupsUpdate)
	t.Run("RoleMenuOptions", testRoleMenuOptionsUpdate)
	t.Run("RoleMenus", testRoleMenusUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsSliceUpdateAll)
	t.Run("RoleGroups", testRoleGroupsSliceUpdateAll)
	t.Run("RoleMenuOptions", testRoleMenuOptionsSliceUpdateAll)
	t.Run("RoleMenus", testRoleMenusSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("RoleCommands", testRoleCommandsUpsert)
	t.Run("RoleGroups", testRoleGroupsUpsert)
	t.Run("RoleMenuOptions", testRoleMenuOptionsUpsert)
	t.Run("RoleMenus", testRoleMenusUpsert)
}