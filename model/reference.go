package model

//
// Relation info: stored independently from the instances they refer to.
//
type References struct {
	classes   ClassMap
	instances InstanceMap
	tables    TableRelations
}

//
// via References.FindByName()
//
type Reference struct {
	*References
	inst *InstanceInfo
}

func NewReferences(classes ClassMap, instances InstanceMap, tables TableRelations) References {
	return References{classes, instances, tables}
}

//
// Add a new set of references for the passed id'd reference.
//
func (this *References) NewInstance(id StringId, class *ClassInfo, name string, long string) *InstanceInfo {
	inst := NewInstanceInfo(id, class, name, long, this)
	this.instances[id] = inst
	return inst
}

//
// Find the set of references used by the named instance.
//
func (this *References) FindByName(name string) (ret Reference, err error) {
	if inst, e := this.instances.FindInstance(name); e != nil {
		err = e
	} else {
		ret = Reference{this, inst}
	}
	return ret, err
}

//
// Test whther the referenced inst is compatible with the passed (class) id.
//
func (this *Reference) CompatibleWith(classId StringId) (okay bool) {
	if class, ok := this.classes[classId]; ok {
		okay = this.inst.CompatibleWith(class)
	}
	return okay
}
