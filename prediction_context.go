package antlr

import (
	"fmt"
	"strconv"
)

type PredictionContext interface {
	Hash() string
	GetParent(int) PredictionContext
	getReturnState(int) int
	equals(PredictionContext) bool
	length() int
	isEmpty() bool
	hasEmptyPath() bool
	String() string
}

type BasePredictionContext struct {
	cachedHashString string
}

func NewBasePredictionContext(cachedHashString string) *BasePredictionContext {

	pc := new(BasePredictionContext)
	pc.cachedHashString = cachedHashString

	return pc
}

// Represents {@code $} in local context prediction, which means wildcard.
// {@code//+x =//}.
// /
const (
	BasePredictionContextEMPTY_RETURN_STATE = 0x7FFFFFFF
)

// Represents {@code $} in an array in full context mode, when {@code $}
// doesn't mean wildcard: {@code $ + x = [$,x]}. Here,
// {@code $} = {@link //EMPTY_RETURN_STATE}.
// /

var BasePredictionContextglobalNodeCount = 1
var BasePredictionContextid = BasePredictionContextglobalNodeCount

// Stores the computed hash code of this {@link BasePredictionContext}. The hash
// code is computed in parts to Match the following reference algorithm.
//
// <pre>
// private int referenceHashCode() {
// int hash = {@link MurmurHash//initialize MurmurHash.initialize}({@link
// //INITIAL_HASH})
//
// for (int i = 0 i &lt {@link //Size()} i++) {
// hash = {@link MurmurHash//update MurmurHash.update}(hash, {@link //GetParent
// GetParent}(i))
// }
//
// for (int i = 0 i &lt {@link //Size()} i++) {
// hash = {@link MurmurHash//update MurmurHash.update}(hash, {@link
// //getReturnState getReturnState}(i))
// }
//
// hash = {@link MurmurHash//finish MurmurHash.finish}(hash, 2// {@link
// //Size()})
// return hash
// }
// </pre>
//

func (this *BasePredictionContext) isEmpty() bool {
	return false
}

func (this *BasePredictionContext) Hash() string {
	return this.cachedHashString
}

func calculateHashString(parent PredictionContext, returnState int) string {
	return parent.String() + strconv.Itoa(returnState)
}

func calculateEmptyHashString() string {
	return ""
}

// Used to cache {@link BasePredictionContext} objects. Its used for the shared
// context cash associated with contexts in DFA states. This cache
// can be used for both lexers and parsers.

type PredictionContextCache struct {
	cache map[PredictionContext]PredictionContext
}

func NewPredictionContextCache() *PredictionContextCache {
	t := new(PredictionContextCache)
	t.cache = make(map[PredictionContext]PredictionContext)
	return t
}

// Add a context to the cache and return it. If the context already exists,
// return that one instead and do not add a Newcontext to the cache.
// Protect shared cache from unsafe thread access.
//
func (this *PredictionContextCache) add(ctx PredictionContext) PredictionContext {
	if ctx == BasePredictionContextEMPTY {
		return BasePredictionContextEMPTY
	}
	var existing = this.cache[ctx]
	if existing != nil {
		return existing
	}
	this.cache[ctx] = ctx
	return ctx
}

func (this *PredictionContextCache) Get(ctx PredictionContext) PredictionContext {
	return this.cache[ctx]
}

func (this *PredictionContextCache) length() int {
	return len(this.cache)
}

type SingletonPredictionContext interface {
	PredictionContext
}

type BaseSingletonPredictionContext struct {
	*BasePredictionContext

	parentCtx   PredictionContext
	returnState int
}

func NewBaseSingletonPredictionContext(parent PredictionContext, returnState int) *BaseSingletonPredictionContext {

	s := new(BaseSingletonPredictionContext)
	s.BasePredictionContext = NewBasePredictionContext("")

	if parent != nil {
		s.cachedHashString = calculateHashString(parent, returnState)
	} else {
		s.cachedHashString = calculateEmptyHashString()
	}

	s.parentCtx = parent
	s.returnState = returnState

	return s
}

func SingletonBasePredictionContextCreate(parent PredictionContext, returnState int) PredictionContext {
	if returnState == BasePredictionContextEMPTY_RETURN_STATE && parent == nil {
		// someone can pass in the bits of an array ctx that mean $
		return BasePredictionContextEMPTY
	} else {
		return NewBaseSingletonPredictionContext(parent, returnState)
	}
}

func (this *BaseSingletonPredictionContext) length() int {
	return 1
}

func (this *BaseSingletonPredictionContext) GetParent(index int) PredictionContext {
	return this.parentCtx
}

func (this *BaseSingletonPredictionContext) getReturnState(index int) int {
	return this.returnState
}

func (this *BaseSingletonPredictionContext) hasEmptyPath() bool {
	return this.returnState == BasePredictionContextEMPTY_RETURN_STATE
}

func (this *BaseSingletonPredictionContext) equals(other PredictionContext) bool {
	if this == other {
		return true
	} else if _, ok := other.(*BaseSingletonPredictionContext); !ok {
		return false
	} else if this.Hash() != other.Hash() {
		return false // can't be same if hash is different
	} else {

		otherP := other.(*BaseSingletonPredictionContext)

		if this.returnState != other.getReturnState(0) {
			return false
		} else if this.parentCtx == nil {
			return otherP.parentCtx == nil
		} else {
			return this.parentCtx.equals(otherP.parentCtx)
		}
	}
}

func (this *BaseSingletonPredictionContext) Hash() string {
	return this.cachedHashString
}

func (this *BaseSingletonPredictionContext) String() string {
	var up string

	if this.parentCtx == nil {
		up = ""
	} else {
		up = this.parentCtx.String()
	}

	if len(up) == 0 {
		if this.returnState == BasePredictionContextEMPTY_RETURN_STATE {
			return "$"
		} else {
			return strconv.Itoa(this.returnState)
		}
	} else {
		return strconv.Itoa(this.returnState) + " " + up
	}
}

var BasePredictionContextEMPTY = NewEmptyPredictionContext()

type EmptyPredictionContext struct {
	*BaseSingletonPredictionContext
}

func NewEmptyPredictionContext() *EmptyPredictionContext {

	p := new(EmptyPredictionContext)

	p.BaseSingletonPredictionContext = NewBaseSingletonPredictionContext(nil, BasePredictionContextEMPTY_RETURN_STATE)

	return p
}

func (this *EmptyPredictionContext) isEmpty() bool {
	return true
}

func (this *EmptyPredictionContext) GetParent(index int) PredictionContext {
	return nil
}

func (this *EmptyPredictionContext) getReturnState(index int) int {
	return this.returnState
}

func (this *EmptyPredictionContext) equals(other PredictionContext) bool {
	return this == other
}

func (this *EmptyPredictionContext) String() string {
	return "$"
}

type ArrayPredictionContext struct {
	*BasePredictionContext

	parents      []PredictionContext
	returnStates []int
}

func NewArrayPredictionContext(parents []PredictionContext, returnStates []int) *ArrayPredictionContext {
	// Parent can be nil only if full ctx mode and we make an array
	// from {@link //EMPTY} and non-empty. We merge {@link //EMPTY} by using
	// nil parent and
	// returnState == {@link //EMPTY_RETURN_STATE}.

	c := new(ArrayPredictionContext)
	c.BasePredictionContext = NewBasePredictionContext("")

	for i := range parents {
		c.cachedHashString += calculateHashString(parents[i], returnStates[i])
	}

	c.parents = parents
	c.returnStates = returnStates

	return c
}

func (c *ArrayPredictionContext) GetReturnStates() []int {
	return c.returnStates
}

func (this *ArrayPredictionContext) hasEmptyPath() bool {
	return this.getReturnState(this.length()-1) == BasePredictionContextEMPTY_RETURN_STATE
}

func (this *ArrayPredictionContext) isEmpty() bool {
	// since EMPTY_RETURN_STATE can only appear in the last position, we
	// don't need to verify that size==1
	return this.returnStates[0] == BasePredictionContextEMPTY_RETURN_STATE
}

func (this *ArrayPredictionContext) length() int {
	return len(this.returnStates)
}

func (this *ArrayPredictionContext) GetParent(index int) PredictionContext {
	return this.parents[index]
}

func (this *ArrayPredictionContext) getReturnState(index int) int {
	return this.returnStates[index]
}

func (this *ArrayPredictionContext) equals(other PredictionContext) bool {
	if _, ok := other.(*ArrayPredictionContext); !ok {
		return false
	} else if this.cachedHashString != other.Hash() {
		return false // can't be same if hash is different
	} else {
		otherP := other.(*ArrayPredictionContext)
		return &this.returnStates == &otherP.returnStates && &this.parents == &otherP.parents
	}
}

func (this *ArrayPredictionContext) String() string {
	if this.isEmpty() {
		return "[]"
	} else {
		var s = "["
		for i := 0; i < len(this.returnStates); i++ {
			if i > 0 {
				s = s + ", "
			}
			if this.returnStates[i] == BasePredictionContextEMPTY_RETURN_STATE {
				s = s + "$"
				continue
			}
			s = s + strconv.Itoa(this.returnStates[i])
			if this.parents[i] != nil {
				s = s + " " + this.parents[i].String()
			} else {
				s = s + "nil"
			}
		}
		return s + "]"
	}
}

// Convert a {@link RuleContext} tree to a {@link BasePredictionContext} graph.
// Return {@link //EMPTY} if {@code outerContext} is empty or nil.
// /
func predictionContextFromRuleContext(a *ATN, outerContext RuleContext) PredictionContext {
	if outerContext == nil {
		outerContext = RuleContextEmpty
	}
	// if we are in RuleContext of start rule, s, then BasePredictionContext
	// is EMPTY. Nobody called us. (if we are empty, return empty)
	if outerContext.GetParent() == nil || outerContext == RuleContextEmpty {
		return BasePredictionContextEMPTY
	}
	// If we have a parent, convert it to a BasePredictionContext graph
	var parent = predictionContextFromRuleContext(a, outerContext.GetParent().(RuleContext))
	var state = a.states[outerContext.GetInvokingState()]
	var transition = state.GetTransitions()[0]

	return SingletonBasePredictionContextCreate(parent, transition.(*RuleTransition).followState.GetStateNumber())
}

func calculateListsHashString(parents []BasePredictionContext, returnStates []int) string {
	var s = ""

	for _, p := range parents {
		s += fmt.Sprint(p)
	}

	for _, r := range returnStates {
		s += fmt.Sprint(r)
	}

	return s
}

func merge(a, b PredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) PredictionContext {
	// share same graph if both same
	if a == b {
		return a
	}

	ac, ok1 := a.(*BaseSingletonPredictionContext)
	bc, ok2 := a.(*BaseSingletonPredictionContext)

	if ok1 && ok2 {
		return mergeSingletons(ac, bc, rootIsWildcard, mergeCache)
	}
	// At least one of a or b is array
	// If one is $ and rootIsWildcard, return $ as// wildcard
	if rootIsWildcard {
		if _, ok := a.(*EmptyPredictionContext); ok {
			return a
		}
		if _, ok := b.(*EmptyPredictionContext); ok {
			return b
		}
	}
	// convert singleton so both are arrays to normalize
	if _, ok := a.(*BaseSingletonPredictionContext); ok {
		a = NewArrayPredictionContext([]PredictionContext{a.GetParent(0)}, []int{a.getReturnState(0)})
	}
	if _, ok := b.(*BaseSingletonPredictionContext); ok {
		b = NewArrayPredictionContext([]PredictionContext{b.GetParent(0)}, []int{b.getReturnState(0)})
	}
	return mergeArrays(a.(*ArrayPredictionContext), b.(*ArrayPredictionContext), rootIsWildcard, mergeCache)
}

//
// Merge two {@link SingletonBasePredictionContext} instances.
//
// <p>Stack tops equal, parents merge is same return left graph.<br>
// <embed src="images/SingletonMerge_SameRootSamePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Same stack top, parents differ merge parents giving array node, then
// remainders of those graphs. A Newroot node is created to point to the
// merged parents.<br>
// <embed src="images/SingletonMerge_SameRootDiffPar.svg"
// type="image/svg+xml"/></p>
//
// <p>Different stack tops pointing to same parent. Make array node for the
// root where both element in the root point to the same (original)
// parent.<br>
// <embed src="images/SingletonMerge_DiffRootSamePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Different stack tops pointing to different parents. Make array node for
// the root where each element points to the corresponding original
// parent.<br>
// <embed src="images/SingletonMerge_DiffRootDiffPar.svg"
// type="image/svg+xml"/></p>
//
// @param a the first {@link SingletonBasePredictionContext}
// @param b the second {@link SingletonBasePredictionContext}
// @param rootIsWildcard {@code true} if this is a local-context merge,
// otherwise false to indicate a full-context merge
// @param mergeCache
// /
func mergeSingletons(a, b *BaseSingletonPredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) PredictionContext {
	if mergeCache != nil {
		var previous = mergeCache.Get(a.Hash(), b.Hash())
		if previous != nil {
			return previous.(PredictionContext)
		}
		previous = mergeCache.Get(b.Hash(), a.Hash())
		if previous != nil {
			return previous.(PredictionContext)
		}
	}

	var rootMerge = mergeRoot(a, b, rootIsWildcard)
	if rootMerge != nil {
		if mergeCache != nil {
			mergeCache.set(a.Hash(), b.Hash(), rootMerge)
		}
		return rootMerge
	}
	if a.returnState == b.returnState {
		var parent = merge(a.parentCtx, b.parentCtx, rootIsWildcard, mergeCache)
		// if parent is same as existing a or b parent or reduced to a parent,
		// return it
		if parent == a.parentCtx {
			return a // ax + bx = ax, if a=b
		}
		if parent == b.parentCtx {
			return b // ax + bx = bx, if a=b
		}
		// else: ax + ay = a'[x,y]
		// merge parents x and y, giving array node with x,y then remainders
		// of those graphs. dup a, a' points at merged array
		// Newjoined parent so create Newsingleton pointing to it, a'
		var spc = SingletonBasePredictionContextCreate(parent, a.returnState)
		if mergeCache != nil {
			mergeCache.set(a.Hash(), b.Hash(), spc)
		}
		return spc
	} else { // a != b payloads differ
		// see if we can collapse parents due to $+x parents if local ctx
		var singleParent PredictionContext = nil
		if a == b || (a.parentCtx != nil && a.parentCtx == b.parentCtx) { // ax +
			// bx =
			// [a,b]x
			singleParent = a.parentCtx
		}
		if singleParent != nil { // parents are same
			// sort payloads and use same parent
			var payloads = []int{a.returnState, b.returnState}
			if a.returnState > b.returnState {
				payloads[0] = b.returnState
				payloads[1] = a.returnState
			}
			var parents = []PredictionContext{singleParent, singleParent}
			var apc = NewArrayPredictionContext(parents, payloads)
			if mergeCache != nil {
				mergeCache.set(a.Hash(), b.Hash(), apc)
			}
			return apc
		}
		// parents differ and can't merge them. Just pack together
		// into array can't merge.
		// ax + by = [ax,by]
		var payloads = []int{a.returnState, b.returnState}
		var parents = []PredictionContext{a.parentCtx, b.parentCtx}
		if a.returnState > b.returnState { // sort by payload
			payloads[0] = b.returnState
			payloads[1] = a.returnState
			parents = []PredictionContext{b.parentCtx, a.parentCtx}
		}
		var a_ = NewArrayPredictionContext(parents, payloads)
		if mergeCache != nil {
			mergeCache.set(a.Hash(), b.Hash(), a_)
		}
		return a_
	}
}

//
// Handle case where at least one of {@code a} or {@code b} is
// {@link //EMPTY}. In the following diagrams, the symbol {@code $} is used
// to represent {@link //EMPTY}.
//
// <h2>Local-Context Merges</h2>
//
// <p>These local-context merge operations are used when {@code rootIsWildcard}
// is true.</p>
//
// <p>{@link //EMPTY} is superset of any graph return {@link //EMPTY}.<br>
// <embed src="images/LocalMerge_EmptyRoot.svg" type="image/svg+xml"/></p>
//
// <p>{@link //EMPTY} and anything is {@code //EMPTY}, so merged parent is
// {@code //EMPTY} return left graph.<br>
// <embed src="images/LocalMerge_EmptyParent.svg" type="image/svg+xml"/></p>
//
// <p>Special case of last merge if local context.<br>
// <embed src="images/LocalMerge_DiffRoots.svg" type="image/svg+xml"/></p>
//
// <h2>Full-Context Merges</h2>
//
// <p>These full-context merge operations are used when {@code rootIsWildcard}
// is false.</p>
//
// <p><embed src="images/FullMerge_EmptyRoots.svg" type="image/svg+xml"/></p>
//
// <p>Must keep all contexts {@link //EMPTY} in array is a special value (and
// nil parent).<br>
// <embed src="images/FullMerge_EmptyRoot.svg" type="image/svg+xml"/></p>
//
// <p><embed src="images/FullMerge_SameRoot.svg" type="image/svg+xml"/></p>
//
// @param a the first {@link SingletonBasePredictionContext}
// @param b the second {@link SingletonBasePredictionContext}
// @param rootIsWildcard {@code true} if this is a local-context merge,
// otherwise false to indicate a full-context merge
// /
func mergeRoot(a, b SingletonPredictionContext, rootIsWildcard bool) PredictionContext {
	if rootIsWildcard {
		if a == BasePredictionContextEMPTY {
			return BasePredictionContextEMPTY // // + b =//
		}
		if b == BasePredictionContextEMPTY {
			return BasePredictionContextEMPTY // a +// =//
		}
	} else {
		if a == BasePredictionContextEMPTY && b == BasePredictionContextEMPTY {
			return BasePredictionContextEMPTY // $ + $ = $
		} else if a == BasePredictionContextEMPTY { // $ + x = [$,x]
			var payloads = []int{b.getReturnState(-1), BasePredictionContextEMPTY_RETURN_STATE}
			var parents = []PredictionContext{b.GetParent(-1), nil}
			return NewArrayPredictionContext(parents, payloads)
		} else if b == BasePredictionContextEMPTY { // x + $ = [$,x] ($ is always first if present)
			var payloads = []int{a.getReturnState(-1), BasePredictionContextEMPTY_RETURN_STATE}
			var parents = []PredictionContext{a.GetParent(-1), nil}
			return NewArrayPredictionContext(parents, payloads)
		}
	}
	return nil
}

//
// Merge two {@link ArrayBasePredictionContext} instances.
//
// <p>Different tops, different parents.<br>
// <embed src="images/ArrayMerge_DiffTopDiffPar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, same parents.<br>
// <embed src="images/ArrayMerge_ShareTopSamePar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, different parents.<br>
// <embed src="images/ArrayMerge_ShareTopDiffPar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, all shared parents.<br>
// <embed src="images/ArrayMerge_ShareTopSharePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Equal tops, merge parents and reduce top to
// {@link SingletonBasePredictionContext}.<br>
// <embed src="images/ArrayMerge_EqualTop.svg" type="image/svg+xml"/></p>
// /
func mergeArrays(a, b *ArrayPredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) PredictionContext {
	if mergeCache != nil {
		var previous = mergeCache.Get(a.Hash(), b.Hash())
		if previous != nil {
			return previous.(PredictionContext)
		}
		previous = mergeCache.Get(b.Hash(), a.Hash())
		if previous != nil {
			return previous.(PredictionContext)
		}
	}
	// merge sorted payloads a + b => M
	var i = 0 // walks a
	var j = 0 // walks b
	var k = 0 // walks target M array

	var mergedReturnStates = make([]int, 0)
	var mergedParents = make([]PredictionContext, 0)
	// walk and merge to yield mergedParents, mergedReturnStates
	for i < len(a.returnStates) && j < len(b.returnStates) {
		var a_parent = a.parents[i]
		var b_parent = b.parents[j]
		if a.returnStates[i] == b.returnStates[j] {
			// same payload (stack tops are equal), must yield merged singleton
			var payload = a.returnStates[i]
			// $+$ = $
			var bothDollars = payload == BasePredictionContextEMPTY_RETURN_STATE && a_parent == nil && b_parent == nil
			var ax_ax = (a_parent != nil && b_parent != nil && a_parent == b_parent) // ax+ax
			// ->
			// ax
			if bothDollars || ax_ax {
				mergedParents[k] = a_parent // choose left
				mergedReturnStates[k] = payload
			} else { // ax+ay -> a'[x,y]
				var mergedParent = merge(a_parent, b_parent, rootIsWildcard, mergeCache)
				mergedParents[k] = mergedParent
				mergedReturnStates[k] = payload
			}
			i += 1 // hop over left one as usual
			j += 1 // but also Skip one in right side since we merge
		} else if a.returnStates[i] < b.returnStates[j] { // copy a[i] to M
			mergedParents[k] = a_parent
			mergedReturnStates[k] = a.returnStates[i]
			i += 1
		} else { // b > a, copy b[j] to M
			mergedParents[k] = b_parent
			mergedReturnStates[k] = b.returnStates[j]
			j += 1
		}
		k += 1
	}
	// copy over any payloads remaining in either array
	if i < len(a.returnStates) {
		for p := i; p < len(a.returnStates); p++ {
			mergedParents[k] = a.parents[p]
			mergedReturnStates[k] = a.returnStates[p]
			k += 1
		}
	} else {
		for p := j; p < len(b.returnStates); p++ {
			mergedParents[k] = b.parents[p]
			mergedReturnStates[k] = b.returnStates[p]
			k += 1
		}
	}
	// trim merged if we combined a few that had same stack tops
	if k < len(mergedParents) { // write index < last position trim
		if k == 1 { // for just one merged element, return singleton top
			var a_ = SingletonBasePredictionContextCreate(mergedParents[0], mergedReturnStates[0])
			if mergeCache != nil {
				mergeCache.set(a.Hash(), b.Hash(), a_)
			}
			return a_
		}
		mergedParents = mergedParents[0:k]
		mergedReturnStates = mergedReturnStates[0:k]
	}

	var M = NewArrayPredictionContext(mergedParents, mergedReturnStates)

	// if we created same array as a or b, return that instead
	// TODO: track whether this is possible above during merge sort for speed
	if M == a {
		if mergeCache != nil {
			mergeCache.set(a.Hash(), b.Hash(), a)
		}
		return a
	}
	if M == b {
		if mergeCache != nil {
			mergeCache.set(a.Hash(), b.Hash(), b)
		}
		return b
	}
	combineCommonParents(mergedParents)

	if mergeCache != nil {
		mergeCache.set(a.Hash(), b.Hash(), M)
	}
	return M
}

//
// Make pass over all <em>M</em> {@code parents} merge any {@code equals()}
// ones.
// /
func combineCommonParents(parents []PredictionContext) {
	var uniqueParents = make(map[PredictionContext]PredictionContext)

	for p := 0; p < len(parents); p++ {
		var parent = parents[p]
		if uniqueParents[parent] == nil {
			uniqueParents[parent] = parent
		}
	}
	for q := 0; q < len(parents); q++ {
		parents[q] = uniqueParents[parents[q]]
	}
}

func getCachedBasePredictionContext(context PredictionContext, contextCache *PredictionContextCache, visited map[PredictionContext]PredictionContext) PredictionContext {

	if context.isEmpty() {
		return context
	}
	var existing = visited[context]
	if existing != nil {
		return existing
	}
	existing = contextCache.Get(context)
	if existing != nil {
		visited[context] = existing
		return existing
	}
	var changed = false
	var parents = make([]PredictionContext, context.length())
	for i := 0; i < len(parents); i++ {
		var parent = getCachedBasePredictionContext(context.GetParent(i), contextCache, visited)
		if changed || parent != context.GetParent(i) {
			if !changed {
				parents = make([]PredictionContext, context.length())
				for j := 0; j < context.length(); j++ {
					parents[j] = context.GetParent(j)
				}
				changed = true
			}
			parents[i] = parent
		}
	}
	if !changed {
		contextCache.add(context)
		visited[context] = context
		return context
	}
	var updated PredictionContext = nil
	if len(parents) == 0 {
		updated = BasePredictionContextEMPTY
	} else if len(parents) == 1 {
		updated = SingletonBasePredictionContextCreate(parents[0], context.getReturnState(0))
	} else {
		updated = NewArrayPredictionContext(parents, context.(*ArrayPredictionContext).GetReturnStates())
	}
	contextCache.add(updated)
	visited[updated] = updated
	visited[context] = updated

	return updated
}

// ter's recursive version of Sam's getAllNodes()
//func getAllContextNodes(context, nodes, visited) {
//	if (nodes == nil) {
//		nodes = []
//		return getAllContextNodes(context, nodes, visited)
//	} else if (visited == nil) {
//		visited = {}
//		return getAllContextNodes(context, nodes, visited)
//	} else {
//		if (context == nil || visited[context] != nil) {
//			return nodes
//		}
//		visited[context] = context
//		nodes.push(context)
//		for i := 0; i < len(context); i++ {
//			getAllContextNodes(context.GetParent(i), nodes, visited)
//		}
//		return nodes
//	}
//}