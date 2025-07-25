/*
 * Copyright 2023 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operations

import (
	"github.com/yorkie-team/yorkie/pkg/document/crdt"
	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/pkg/resource"
)

// TreeStyle represents the style operation of the tree.
type TreeStyle struct {
	// parentCreatedAt is the creation time of the Text that executes Style.
	parentCreatedAt *time.Ticket

	// fromPos represents the start point of the editing range.
	from *crdt.TreePos

	// toPos represents the end point of the editing range.
	to *crdt.TreePos

	// attributes represents the tree style to be added.
	attributes map[string]string

	// attributesToRemove represents the tree style to be removed.
	attributesToRemove []string

	// executedAt is the time the operation was executed.
	executedAt *time.Ticket
}

// NewTreeStyle creates a new instance of TreeStyle.
func NewTreeStyle(
	parentCreatedAt *time.Ticket,
	from *crdt.TreePos,
	to *crdt.TreePos,
	attributes map[string]string,
	executedAt *time.Ticket,
) *TreeStyle {
	return &TreeStyle{
		parentCreatedAt:    parentCreatedAt,
		from:               from,
		to:                 to,
		attributes:         attributes,
		attributesToRemove: []string{},
		executedAt:         executedAt,
	}
}

// NewTreeStyleRemove creates a new instance of TreeStyle.
func NewTreeStyleRemove(
	parentCreatedAt *time.Ticket,
	from *crdt.TreePos,
	to *crdt.TreePos,
	attributesToRemove []string,
	executedAt *time.Ticket,
) *TreeStyle {
	return &TreeStyle{
		parentCreatedAt:    parentCreatedAt,
		from:               from,
		to:                 to,
		attributes:         map[string]string{},
		attributesToRemove: attributesToRemove,
		executedAt:         executedAt,
	}
}

// Execute executes this operation on the given `CRDTRoot`.
func (e *TreeStyle) Execute(root *crdt.Root, versionVector time.VersionVector) error {
	parent := root.FindByCreatedAt(e.parentCreatedAt)
	obj, ok := parent.(*crdt.Tree)
	if !ok {
		return ErrNotApplicableDataType
	}

	var pairs []crdt.GCPair
	var diff resource.DataSize
	var err error
	if len(e.attributes) > 0 {
		pairs, diff, err = obj.Style(e.from, e.to, e.attributes, e.executedAt, versionVector)
		if err != nil {
			return err
		}
	} else {
		pairs, diff, err = obj.RemoveStyle(e.from, e.to, e.attributesToRemove, e.executedAt, versionVector)
		if err != nil {
			return err
		}
	}

	for _, pair := range pairs {
		root.RegisterGCPair(pair)
		root.AdjustDiffForGCPair(&diff, pair)
	}

	root.Acc(diff)

	return nil
}

// FromPos returns the start point of the editing range.
func (e *TreeStyle) FromPos() *crdt.TreePos {
	return e.from
}

// ToPos returns the end point of the editing range.
func (e *TreeStyle) ToPos() *crdt.TreePos {
	return e.to
}

// ExecutedAt returns execution time of this operation.
func (e *TreeStyle) ExecutedAt() *time.Ticket {
	return e.executedAt
}

// SetActor sets the given actor to this operation.
func (e *TreeStyle) SetActor(actorID time.ActorID) {
	e.executedAt = e.executedAt.SetActorID(actorID)
}

// ParentCreatedAt returns the creation time of the Text.
func (e *TreeStyle) ParentCreatedAt() *time.Ticket {
	return e.parentCreatedAt
}

// Attributes returns the content of Style.
func (e *TreeStyle) Attributes() map[string]string {
	return e.attributes
}

// AttributesToRemove returns the content of Style.
func (e *TreeStyle) AttributesToRemove() []string {
	return e.attributesToRemove
}
