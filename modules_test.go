package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// REGISTER EMPTY MODULE

func TestModules_000(t *testing.T) {
	// Should panic when there is no new function
	defer func() {
		if err := recover(); err != nil {
			t.Log("Received error:", err)
		}
	}()
	gopi.RegisterModule(gopi.Module{})
	t.Error("Expected failure without NewFunction")
}

func TestModules_001(t *testing.T) {
	// Should panic when there is no module name
	defer func() {
		if err := recover(); err != nil {
			t.Log("Received error:", err)
		}
	}()
	gopi.RegisterModule(gopi.Module{
		New: EmptyModuleNewFunction,
	})
	t.Error("Expected failure without module name")
}

func TestModules_002(t *testing.T) {
	// Should panic when the module name is a reserved word
	defer func() {
		if err := recover(); err != nil {
			t.Log("Received error:", err)
		}
	}()
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Name: "gpio",
	})
	t.Error("Expected failure with reserved module word")
}

func TestModules_003(t *testing.T) {
	// Should panic when module is registered twice
	defer func() {
		if err := recover(); err != nil {
			t.Log("Received error:", err)
		}
	}()
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Type: gopi.MODULE_TYPE_GPIO,
	})
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Type: gopi.MODULE_TYPE_GPIO,
	})
	t.Error("Expected failure with double registration")
}

func TestModules_004(t *testing.T) {
	// Should panic when module is registered twice
	defer func() {
		if err := recover(); err != nil {
			t.Log("Received error:", err)
		}
	}()
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Name: "test",
	})
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Name: "test",
	})
	t.Error("Expected failure with double registration")
}

func TestModules_005(t *testing.T) {
	// Should return error when module dependency is not satisfied
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test2",
		Requires: []string{"test3"},
	})

	if _, err := gopi.ModuleWithDependencies("test2"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with unmet dependency")
	}
}

func TestModules_006(t *testing.T) {
	// Should return error when module dependency is on itself
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test3",
		Requires: []string{"test3"},
	})

	if _, err := gopi.ModuleWithDependencies("test3"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with dependency on self")
	}
}

func TestModules_007(t *testing.T) {
	// Should return error when module dependency is circular
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test4",
		Requires: []string{"test5"},
	})
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test5",
		Requires: []string{"test4"},
	})

	if _, err := gopi.ModuleWithDependencies("test4"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with circular dependencies (test4)")
	}
	if _, err := gopi.ModuleWithDependencies("test5"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with circular dependencies (test5)")
	}
}

func TestModules_008(t *testing.T) {
	// Should return error when module dependency is circular
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test6",
		Requires: []string{"test7"},
	})
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test7",
		Requires: []string{"test8"},
	})
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test8",
		Requires: []string{"test6"},
	})

	if _, err := gopi.ModuleWithDependencies("test6"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with circular dependencies (test6)")
	}
	if _, err := gopi.ModuleWithDependencies("test7"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with circular dependencies (test7)")
	}
	if _, err := gopi.ModuleWithDependencies("test8"); err != nil {
		t.Log("Received error:", err)
	} else {
		t.Error("Expected failure with circular dependencies (test8)")
	}
}

func TestModules_009(t *testing.T) {
	// Check for correct order of modules returned
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test9",
		Requires: []string{"test10"},
	})
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test10",
		Requires: []string{"test11"},
	})
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Name: "test11",
	})

	if modules, err := gopi.ModuleWithDependencies("test9"); err != nil {
		t.Error("Received error:", err)
	} else {
		if len(modules) != 3 {
			t.Fatalf("Expected three modules to be returned, got %v", modules)
		}
		if modules[0].Name != "test11" {
			t.Errorf("Expected test11 as the first module, got %v", modules)
		}
		if modules[1].Name != "test10" {
			t.Errorf("Expected test10 as the second module, got %v", modules)
		}
		if modules[2].Name != "test9" {
			t.Errorf("Expected test9 as the third module, got %v", modules)
		}
	}
}

func TestModules_010(t *testing.T) {
	// Check for correct order of modules returned, and two modules
	// depending on a third one
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test12",
		Requires: []string{"test13", "test14"},
	})
	gopi.RegisterModule(gopi.Module{
		New:      EmptyModuleNewFunction,
		Name:     "test13",
		Requires: []string{"test14"},
	})
	gopi.RegisterModule(gopi.Module{
		New:  EmptyModuleNewFunction,
		Name: "test14",
	})

	if modules, err := gopi.ModuleWithDependencies("test12"); err != nil {
		t.Error("Received error:", err)
	} else {
		if len(modules) != 3 {
			t.Fatalf("Expected three modules to be returned, got %v", modules)
		}
		if modules[0].Name != "test14" {
			t.Errorf("Expected test14 as the first module, got %v", modules)
		}
		if modules[1].Name != "test13" {
			t.Errorf("Expected test13 as the second module, got %v", modules)
		}
		if modules[2].Name != "test12" {
			t.Errorf("Expected test12 as the third module, got %v", modules)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// MOCK NEW FUNCTION

func EmptyModuleNewFunction(app *gopi.AppInstance) (gopi.Driver, error) {
	return nil, nil
}
