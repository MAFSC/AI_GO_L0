package bitcoin
type Adapter struct{ mode string }
func NewAdapter(_,_,_,_,mode string)*Adapter{ return &Adapter{mode:mode} }
func (a *Adapter) Mode() string { return a.mode }
