// brigolagecms/pkg/storage/memory/tables.go
//
// Copyright (c) 2020, Michael D Henderson.
// Copyright (c) 2002-2009 Kineticode, Inc. and others.
// Copyright (c) 2001-2002 About.com.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package memory

import "time"

//  -----------------------------------------------------------------------------
//  Table: atType
//
//  Description:
//    The table that holds the information for a given asset type.
//    Holds name and description information and is references by elementType rows.
//
func (s *Storage) newAtType() *atType {
	o := atType{
		id:           s.sequences.atType.next(), // seqAtType.NextVal(),
		topLevel:     false,
		paginated:    false,
		fixedURL:     false,
		relatedStory: false,
		relatedMedia: false,
		media:        false,
		active:       true,
	}
	return &o
}

func (s *Storage) getAtType(id int) *atType {
	if o, ok := s.tables.atType[id]; !ok {
		return nil
	} else {
		return o
	}
}

//  -----------------------------------------------------------------------------
//  TABLE: atTypeMember
//
func (s *Storage) newAtTypeMember() *atTypeMember {
	o := atTypeMember{
		id: s.sequences.atTypeMember.next(), // seqAtTypeMember.NextVal(),
	}
	return &o
}

func (s *Storage) getAtTypeMember(id int) *atTypeMember {
	if o, ok := s.tables.atTypeMember[id]; !ok {
		return nil
	} else {
		return o
	}
}

//  -----------------------------------------------------------------------------
//  Table media
//
//  Description: The Media table this houses the data for a given media asset
//                               and its related asset_version_data
//
func (s *Storage) newMedia() *media {
	o := media{
		id:            s.sequences.media.next(), // seqMedia.NextVal(),
		publishStatus: false,
		active:        true,
	}
	// check aliasID ( aliasID  !  =  id )
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaInstance
//
//  Description: An instance of a media object
//
func (s *Storage) newMediaInstance() *mediaInstance {
	o := mediaInstance{
		id:         s.sequences.mediaInstance.next(), // seqMediaInstance.NextVal(),
		priority:   3,
		coverDate:  time.Now(),
		checkedOut: false,
	}
	// check priority ( priority  between  1  and  5 )
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaURI
//
//  Description: Tracks all URIs for stories.
//
func (s *Storage) newMediaURI() *mediaURI {
	o := mediaURI{
		id: s.sequences.mediaURI.next(), // seqMediaURI.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaOutputChannel
//
//  Description: Mapping Table between stories and output channels.
//
func (s *Storage) newMediaOutputChannel() *mediaOutputChannel {
	o := mediaOutputChannel{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaFields
//
//  Description: A mapping table between Media classes and functions that
//                               Will be run against uploaded files
func (s *Storage) newMediaFields() *mediaFields {
	o := mediaFields{
		id:     s.sequences.mediaFields.next(), // seqMediaFields.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaContributor
//
//  Description: mapping tables between media instances and contributors
//
func (s *Storage) newMediaContributor() *mediaContributor {
	o := mediaContributor{
		id: s.sequences.mediaContributor.next(), // seqMediaContributor.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaMember
//
//  Description: The link between media objects and member objects
//
func (s *Storage) newMediaMember() *mediaMember {
	o := mediaMember{
		id: s.sequences.mediaMember.next(), // seqMediaMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: story
//
//  Description: The story properties. Versioning info might get added here and
//               the rights info might get removed. It is also possible that
//               the asset type field will need a cascading delete.
//
func (s *Storage) newStory() *story {
	o := story{
		id:            s.sequences.story.next(), // seqStory.NextVal(),
		publishStatus: false,
		active:        true,
	}
	// check aliasID ( aliasID  !  =  id )
	return &o
}

//  ----------------------------------------------------------------------------
//  Table storyInstance
//
//  Description:  An instance of a story
//
func (s *Storage) newStoryInstance() *storyInstance {
	o := storyInstance{
		id:         s.sequences.storyInstance.next(), // seqStoryInstance.NextVal(),
		priority:   3,
		coverDate:  time.Now(),
		checkedOut: false,
	}
	// check priority ( priority  between  1  and  5 )
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyURI
//
//  Description: Tracks all URIs for stories.
//
func (s *Storage) newStoryURI() *storyURI {
	o := storyURI{
		id: s.sequences.storyURI.next(), // seqStoryURI.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyOutputChannel
//
//  Description: Mapping Table between stories and output channels.
//
func (s *Storage) newStoryOutputChannel(storyInstanceID, outputChannelID int) *storyOutputChannel {
	o := storyOutputChannel{
		storyInstanceID: storyInstanceID,
		outputChannelID: outputChannelID,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyCategory
//
//  Description: Mapping Table between Stories and categories
//
func (s *Storage) newStoryCategory(storyInstanceID, categoryID int) *storyCategory {
	o := storyCategory{
		id:              s.sequences.storyCategory.next(), // seqStoryCategory.NextVal(),
		storyInstanceID: storyInstanceID,
		categoryID:      categoryID,
		main:            false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyContributor
//
//  Description: mapping tables between story instances and contributors
//
func (s *Storage) newStoryContributor() *storyContributor {
	o := storyContributor{
		id: s.sequences.storyContributor.next(), // seqStoryContributor.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table template
//
//  Description: The table that holds all the template info
//
func (s *Storage) newTemplate() *template {
	o := template{
		id:           s.sequences.template.next(), // seqTemplate.NextVal(),
		tplateType:   1,
		deployStatus: false,
		active:       true,
	}
	// check tplateType ( tplateType  in  1  ,  2  ,  3 )
	return &o
}

//  -----------------------------------------------------------------------------
//  Table templateInstance
//
//  Description:  An versioned instance of a template asset
//
func (s *Storage) newTemplateInstance() *templateInstance {
	o := templateInstance{
		id:         s.sequences.templateInstance.next(), // seqTemplateInstance.NextVal(),
		priority:   3,
		checkedOut: false,
	}
	// check priority ( priority  between  1  and  5 )
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: templateMember
//
//  Description: The link between template objects and member objects
//
func (s *Storage) newTemplateMember() *templateMember {
	o := templateMember{
		id: s.sequences.templateMember.next(), // seqTemplateMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: category
//
//  Description: The category table
//
func (s *Storage) newCategory() *category {
	o := category{
		id:     s.sequences.category.next(), // seqCategory.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: categoryMember
//
//  Description: The link between desk objects and member objects
//
func (s *Storage) newCategoryMember() *categoryMember {
	o := categoryMember{
		id: s.sequences.categoryMember.next(), // seqCategoryMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategory
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its category ID and an attribute name.
//
func (s *Storage) newAttrCategory() *attrCategory {
	o := attrCategory{
		id:     s.sequences.attrCategory.next(), // seqAttrCategory.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategoryVal
//
//  Description: A table to hold attribute values.
//
func (s *Storage) newAttrCategoryVal() *attrCategoryVal {
	o := attrCategoryVal{
		id:     s.sequences.attrCategoryVal.next(), // seqAttrCategoryVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategoryMeta
//
//  Description: A table to represent metadata on types of attributes.
//
func (s *Storage) newAttrCategoryMeta() *attrCategoryMeta {
	o := attrCategoryMeta{
		id:     s.sequences.attrCategoryMeta.next(), // seqAttrCategoryMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: contact
//
func (s *Storage) newContact() *contact {
	o := contact{
		id:        s.sequences.contact.next(), // seqContact.NextVal(),
		active:    true,
		alertable: false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: contactValue
//
func (s *Storage) newContactValue() *contactValue {
	o := contactValue{
		id:     s.sequences.contactValue.next(), // seqContactValue.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: person__contact
//
func (s *Storage) newPersonContactValue() *personContactValue {
	o := personContactValue{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyElement
//
//  Description: Holds the properties of a container element. Note that
//               elements can hold either other container elements or field
//               elements, not both.
//
func (s *Storage) newStoryElement() *storyElement {
	o := storyElement{
		id:        s.sequences.storyElement.next(), // seqStoryElement.NextVal(),
		displayed: false,
		active:    true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaElement
//
//  Description: Holds the properties of a media container element.
//
func (s *Storage) newMediaElement() *mediaElement {
	o := mediaElement{
		id:        s.sequences.mediaElement.next(), // seqMediaElement.NextVal(),
		displayed: false,
		active:    true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyField
//
//  Description: Story Field elements are story specific mappings to the
//               Bric::Biz::Element::Field class. They link to the story that
//               this element is a part of, the attribute id of the data that
//               is contained with in, and it's parent's id (a storyElement
//               row). Place is it's order and active is it's active state.
//
func (s *Storage) newStoryField() *storyField {
	o := storyField{
		id:      s.sequences.storyField.next(), // seqStoryField.NextVal(),
		holdVal: false,
		active:  true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaField
//
//  Description: Media Field elements are media specific mappings to the
//               Bric::Biz::Element::Field class. They link to the media that
//               this element is a part of, the attribute id of the data that
//               is contained with in, and it's parent's id (a mediaElement
//               row). Place is it's order and active is it's active state.
//
func (s *Storage) newMediaField() *mediaField {
	o := mediaField{
		id:      s.sequences.mediaField.next(), // seqMediaField.NextVal(),
		holdVal: false,
		active:  true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: elementType
//
//  Description: The table that holds the information for a given asset type.
//               Holds name and description information and is references by
//               element_contaner and fieldType rows.
//
func (s *Storage) newElementType() *elementType {
	o := elementType{
		id:           s.sequences.elementType.next(), // seqElementType.NextVal(),
		topLevel:     false,
		paginated:    false,
		fixedURI:     false,
		relatedStory: false,
		relatedMedia: false,
		displayed:    false,
		media:        false,
		active:       true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: subelementType
//
//  Description: A table that manages element type parent/child relationships.
//
func (s *Storage) newSubelementType() *subelementType {
	o := subelementType{
		id:            s.sequences.subelementType.next(), // seqSubelementType.NextVal(),
		place:         1,
		minOccurrence: 0,
		maxOccurrence: 0,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: elementTypeSite
//
//  Description: A table that maps
//
func (s *Storage) newElementTypeSite() *elementTypeSite {
	o := elementTypeSite{
		id:     s.sequences.elementTypeSite.next(), // seqElementTypeSite.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: elementTypeOutputChannel
//
//  Description: Holds a reference to the asset type table, the output channel
//               table and an active flag
//
func (s *Storage) newElementTypeOutputChannel() *elementTypeOutputChannel {
	o := elementTypeOutputChannel{
		id:      s.sequences.elementTypeOutputChannel.next(), // seqElementTypeOutputChannel.NextVal(),
		enabled: true,
		active:  true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: elementTypeMember
//
//  Description: The link between element objects and member objects
//
func (s *Storage) newElementTypeMember() *elementTypeMember {
	o := elementTypeMember{
		id: s.sequences.elementTypeMember.next(), // seqElementTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its element ID and an attribute name.
//
func (s *Storage) newAttrElementType() *attrElementType {
	o := attrElementType{
		id:     s.sequences.attrElementType.next(), // seqAttrElementType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element_val
//
//  Description: A table to hold attribute values.
//
func (s *Storage) newAttrElementTypeVal() *attrElementTypeVal {
	o := attrElementTypeVal{
		id:     s.sequences.attrElementTypeVal.next(), // seqAttrElementTypeVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element_meta
//
//  Description: A table to represent metadata on types of attributes.
//
func (s *Storage) newAttrElementTypeMeta() *attrElementTypeMeta {
	o := attrElementTypeMeta{
		id:     s.sequences.attrElementTypeMeta.next(), // seqAttrElementTypeMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: fieldType
//
//  Description: This is the table that contains the name and rules for fields
//          types. It contains references to the elementType table. The place
//          column represents the order that this is to be represented with in
//          its container.
//
func (s *Storage) newFieldType() *fieldType {
	o := fieldType{
		id:            s.sequences.fieldType.next(), // seqFieldType.NextVal(),
		minOccurrence: 0,
		maxOccurrence: 0,
		autopopulated: false,
		maxLength:     0,
		sqlType:       "short",
		widgetType:    "text",
		multiple:      false,
		active:        true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrFieldType
//
//  Description: A table to represent types of attributes. A type is defined by
//               its subsystem, its fieldType ID and an attribute name.
//
func (s *Storage) newAttrFieldType() *attrFieldType {
	o := attrFieldType{
		id:     s.sequences.attrFieldType.next(), // seqAttrFieldType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrFieldTypeVal
//
//  Description: A table to hold attribute values.
//
func (s *Storage) newAttrFieldTypeVal() *attrFieldTypeVal {
	o := attrFieldTypeVal{
		id:     s.sequences.attrFieldTypeVal.next(), // seqAttrFieldTypeVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrFieldTypeMeta
//
//  Description: A table to represent metadata on types of attributes.
//
func (s *Storage) newAttrFieldTypeMeta() *attrFieldTypeMeta {
	o := attrFieldTypeMeta{
		id:     s.sequences.attrFieldTypeMeta.next(), // seqAttrFieldTypeMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: KEYWORD
//
//  Description: The main keyword table.
//
func (s *Storage) newKeyword() *keyword {
	o := keyword{
		id:     s.sequences.keyword.next(), // seqKeyword.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: storyKeyword
//
//  Description: The link between stories and keywords
//
func (s *Storage) newStoryKeyword() *storyKeyword {
	o := storyKeyword{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaKeyword
//
//  Description: The link between media and keywords
//
func (s *Storage) newMediaKeyword() *mediaKeyword {
	o := mediaKeyword{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: categoryKeyword
//
//  Description: The link between categories and keywords
//
func (s *Storage) newCategoryKeyword() *categoryKeyword {
	o := categoryKeyword{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: keywordMember
//
func (s *Storage) newKeywordMember() *keywordMember {
	o := keywordMember{
		id: s.sequences.keywordMember.next(), // seqKeywordMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: org
//
func (s *Storage) newOrg() *org {
	o := org{
		id:       s.sequences.org.next(), // seqOrg.NextVal(),
		personal: false,
		active:   true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: addr
//
func (s *Storage) newAddr() *addr {
	o := addr{
		id:     s.sequences.addr.next(), // seqAddr.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: addrPartType
//
func (s *Storage) newAddrPartType() *addrPartType {
	o := addrPartType{
		id:     s.sequences.addrPartType.next(), // seqAddrPartType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: addrPart
//
func (s *Storage) newAddrPart() *addrPart {
	o := addrPart{
		id: s.sequences.addrPart.next(), // seqAddrPart.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: personOrgAddr
//
func (s *Storage) newPersonOrgAddr() *personOrgAddr {
	o := personOrgAddr{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: personOrg
//
func (s *Storage) newPersonOrg() *personOrg {
	o := personOrg{
		id:     s.sequences.personOrg.next(), // seqPersonOrg.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: source
//
func (s *Storage) newSource() *source {
	o := source{
		id:     s.sequences.source.next(), // seqSource.NextVal(),
		expire: 0,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table outputChannel
//
//  Description: Holds info on the various output channels and is referenced
//                  by templates and elements
//
func (s *Storage) newOutputChannel() *outputChannel {
	o := outputChannel{
		id:      s.sequences.outputChannel.next(), // seqOutputChannel.NextVal(),
		uriCase: 1,
		useSlug: false,
		burner:  1,
		active:  true,
	}
	// check uriCase ( uriCase  in  1  ,  2  ,  3 )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: outputChannelInclude
//
func (s *Storage) newOutputChannelInclude() *outputChannelInclude {
	o := outputChannelInclude{
		id: s.sequences.outputChannelInclude.next(), // seqOutputChannelInclude.NextVal(),
	}
	// check includeOcID ( includeOcID  <  >  outputChannelID )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: outputChannelMember
//
func (s *Storage) newOutputChannelMember() *outputChannelMember {
	o := outputChannelMember{
		id: s.sequences.outputChannelMember.next(), // seqOutputChannelMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: person
//
func (s *Storage) newPerson() *person {
	o := person{
		id:     s.sequences.person.next(), // seqPerson.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: personMember
//
func (s *Storage) newPersonMember() *personMember {
	o := personMember{
		id: s.sequences.personMember.next(), // seqPersonMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: usr
//
func (s *Storage) newUsr() *usr {
	o := usr{
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: siteMember
//
func (s *Storage) newSiteMember() *siteMember {
	o := siteMember{
		id: s.sequences.siteMember.next(), // seqSiteMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: workflow
//
//  Description: The main workflow table.
//
func (s *Storage) newWorkflow() *workflow {
	o := workflow{
		id:     s.sequences.workflow.next(), // seqWorkflow.NextVal(),
		ttypee: 1,
		active: true,
	}
	// check type ( type  in  1  ,  2  ,  3 )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: workflowMember
//
func (s *Storage) newWorkflowMember() *workflowMember {
	o := workflowMember{
		id: s.sequences.workflowMember.next(), // seqWorkflowMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: desk
//
//  Description: Represents a desk in the workflow
//
func (s *Storage) newDesk() *desk {
	o := desk{
		id:      s.sequences.desk.next(), // seqDesk.NextVal(),
		publish: false,
		active:  true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: deskMember
//
//  Description: The link between desk objects and member objects
func (s *Storage) newDeskMember() *deskMember {
	o := deskMember{
		id: s.sequences.deskMember.next(), // seqDeskMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: action
//
func (s *Storage) newAction() *action {
	o := action{
		id:     s.sequences.action.next(), // seqAction.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: attrAction
//
func (s *Storage) newAttrAction() *attrAction {
	o := attrAction{
		id:     s.sequences.attrAction.next(), // seqAttrAction.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: attrActionVal
//
func (s *Storage) newAttrActionVal() *attrActionVal {
	o := attrActionVal{
		id:     s.sequences.attrActionVal.next(), // seqAttrActionVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: attrActionMeta
//
func (s *Storage) newAttrActionMeta() *attrActionMeta {
	o := attrActionMeta{
		id:     s.sequences.attrActionMeta.next(), // seqAttrActionMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: actionType
//
func (s *Storage) newActionType() *actionType {
	o := actionType{
		id:     s.sequences.actionType.next(), // seqActionType.NextVal(),
		active: false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: actionTypeMediaType
//
func (s *Storage) newActionTypeMediaType() *actionTypeMediaType {
	o := actionTypeMediaType{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: mediaResource
//
func (s *Storage) newMediaResource() *mediaResource {
	o := mediaResource{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: resource
//
func (s *Storage) newResource() *resource {
	o := resource{
		id: s.sequences.resource.next(), // seqResource.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: storyResource
//
func (s *Storage) newStoryResource() *storyResource {
	o := storyResource{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: server
//
func (s *Storage) newServer() *server {
	o := server{
		id:     s.sequences.server.next(), // seqServer.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: serverType
//
func (s *Storage) newServerType() *serverType {
	o := serverType{
		id:       s.sequences.serverType.next(), // seqServerType.NextVal(),
		copyable: false,
		publish:  true,
		preview:  false,
		active:   true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: serverType__output__channel
//
func (s *Storage) newServerTypeOutputChannel() *serverTypeOutputChannel {
	o := serverTypeOutputChannel{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: destMember
//
func (s *Storage) newDestMember() *destMember {
	o := destMember{
		id: s.sequences.destMember.next(), // seqDestMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alert
//
func (s *Storage) newAlert() *alert {
	o := alert{
		id:        s.sequences.alert.next(), // seqAlert.NextVal(),
		timestamp: time.Now(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertType
//
func (s *Storage) newAlertType() *alertType {
	o := alertType{
		id:     s.sequences.alertType.next(), // seqAlertType.NextVal(),
		active: true,
		del:    false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertTypeGrpContact
//
func (s *Storage) newAlertTypeGrpContact() *alertTypeGrpContact {
	o := alertTypeGrpContact{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertTypeUsrContact
//
func (s *Storage) newAlertTypeUsrContact() *alertTypeUsrContact {
	o := alertTypeUsrContact{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertTypeRule
//
func (s *Storage) newAlertTypeRule() *alertTypeRule {
	o := alertTypeRule{
		id: s.sequences.alertTypeRule.next(), // seqAlertTypeRule.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alerted
//
func (s *Storage) newAlerted() *alerted {
	o := alerted{
		id: s.sequences.alerted.next(), // seqAlerted.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertedContactValue
//
func (s *Storage) newAlertedContactValue() *alertedContactValue {
	o := alertedContactValue{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: class
//         For keeping track of Perl classes.
//
func (s *Storage) newClass() *class {
	o := class{
		id:          s.sequences.class.next(), // seqClass.NextVal(),
		distributor: false,
	}
	// check keyName ( lower  keyName  =  keyName )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: event
//
func (s *Storage) newEvent() *event {
	o := event{
		id:        s.sequences.event.next(), // seqEvent.NextVal(),
		timestamp: time.Now(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: eventAttr
//
func (s *Storage) newEventAttr() *eventAttr {
	o := eventAttr{
		id: s.sequences.eventAttr.next(), // seqEventAttr.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: eventType
//
func (s *Storage) newEventType() *eventType {
	o := eventType{
		id:     s.sequences.eventType.next(), // seqEventType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: eventTypeAttr
//
func (s *Storage) newEventTypeAttr() *eventTypeAttr {
	o := eventTypeAttr{
		id: s.sequences.eventTypeAttr.next(), // seqEventTypeAttr.NextVal(),
	}
	return &o
}

//  ----------------------------------------------------------------------------
//  Table grp
//
//  Description: The grp table   Contains the name and description of the
//                  group and its parent if it has one
//
func (s *Storage) newGrp() *grp {
	o := grp{
		id:        s.sequences.grp.next(), // seqGrp.NextVal(),
		secret:    true,
		permanent: false,
		active:    true,
	}
	// check parentID ( parentID  <  >  id )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: alertTypeMember
//
func (s *Storage) newAlertTypeMember() *alertTypeMember {
	o := alertTypeMember{
		id: s.sequences.alertTypeMember.next(), // seqAlertTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: contribTypeMember
//
func (s *Storage) newContribTypeMember() *contribTypeMember {
	o := contribTypeMember{
		id: s.sequences.contribTypeMember.next(), // seqContribTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: orgMember
//
func (s *Storage) newOrgMember() *orgMember {
	o := orgMember{
		id: s.sequences.orgMember.next(), // seqOrgMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: grpMember
//
//  Description: The link between stroy objects and member objects
//
func (s *Storage) newGrpMember() *grpMember {
	o := grpMember{
		id: s.sequences.grpMember.next(), // seqGrpMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: member
//
//  Description:    The table that creates a member of a group.   The obj_member
//  table then links the objects to the member table
//
func (s *Storage) newMember() *member {
	o := member{
		id:     s.sequences.member.next(), // seqMember.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: storyMember
//
//  Description: The link between story objects and member objects
//
func (s *Storage) newStoryMember() *storyMember {
	o := storyMember{
		id: s.sequences.storyMember.next(), // seqStoryMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrMember
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its member ID and an attribute name.
//
func (s *Storage) newAttrMember() *attrMember {
	o := attrMember{
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrMemberVal
//  Description: A table to hold attribute values.
//
func (s *Storage) newAttrMemberVal() *attrMemberVal {
	o := attrMemberVal{
		serial: false,
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrMemberMeta
//  Description: A table to represent metadata on types of attributes.
//
func (s *Storage) newAttrMemberMeta() *attrMemberMeta {
	o := attrMemberMeta{
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: sourceMember
//
func (s *Storage) newSourceMember() *sourceMember {
	o := sourceMember{
		id: s.sequences.sourceMember.next(), // seqSourceMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: userMember
//
func (s *Storage) newUserMember() *userMember {
	o := userMember{
		id: s.sequences.userMember.next(), // seqUserMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrGrp
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its grp ID and an attribute name.
//
func (s *Storage) newAttrGrp() *attrGrp {
	o := attrGrp{
		id:     s.sequences.attrGrp.next(), // seqAttrGrp.NextVal(),
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrGrpVal
//
//  Description: A table to hold attribute values.
//
func (s *Storage) newAttrGrpVal() *attrGrpVal {
	o := attrGrpVal{
		id:     s.sequences.attrGrpVal.next(), // seqAttrGrpVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrGrpMeta
//
//  Description: A table to represent metadata on types of attributes.
//
func (s *Storage) newAttrGrpMeta() *attrGrpMeta {
	o := attrGrpMeta{
		id:     s.sequences.attrGrpMeta.next(), // seqAttrGrpMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: job
//
func (s *Storage) newJob() *job {
	o := job{
		id:        s.sequences.job.next(), // seqJob.NextVal(),
		schedTime: time.Now(),
		priority:  3,
		expire:    false,
		failed:    false,
		tries:     0,
		executing: false,
	}
	// check priority ( priority  between  1  and  5 )
	// check tries ( tries  between  0  and  10 )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: jobResource
//
func (s *Storage) newJobResource() *jobResource {
	o := jobResource{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: jobServerType
//
func (s *Storage) newJobServerType() *jobServerType {
	o := jobServerType{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: jobMember
//
func (s *Storage) newJobMember() *jobMember {
	o := jobMember{
		id: s.sequences.jobMember.next(), // seqJobMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: mediaType
//
func (s *Storage) newMediaType() *mediaType {
	o := mediaType{
		id:     s.sequences.mediaType.next(), // seqMediaType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: mediaTypeExt
//
func (s *Storage) newMediaTypeExt() *mediaTypeExt {
	o := mediaTypeExt{
		id: s.sequences.mediaTypeExt.next(), // seqMediaTypeExt.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: mediaTypeMember
//
func (s *Storage) newMediaTypeMember() *mediaTypeMember {
	o := mediaTypeMember{
		id: s.sequences.mediaTypeMember.next(), // seqMediaTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: pref
//         Global preferences.
//
func (s *Storage) newPref() *pref {
	o := pref{
		id:              s.sequences.pref.next(), // seqPref.NextVal(),
		manual:          false,
		canBeOverridden: false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: usrPref
//         Preferences overridden by a specific usr.
//
func (s *Storage) newUsrPref() *usrPref {
	o := usrPref{
		id: s.sequences.usrPref.next(), // seqUsrPref.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: pref
//         Preference options.
//
func (s *Storage) newPrefOpt() *prefOpt {
	o := prefOpt{}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: prefMember
//
func (s *Storage) newPrefMember() *prefMember {
	o := prefMember{
		id: s.sequences.prefMember.next(), // seqPrefMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: grpPriv
//         Privileges granted to user groups.
//
func (s *Storage) newGrpPriv() *grpPriv {
	o := grpPriv{
		id:    s.sequences.priv.next(), // seqPriv.NextVal(),
		mtime: time.Now(),
	}
	// check value ( value  between  1  and  255 )
	return &o
}

//  -----------------------------------------------------------------------------
//  TABLE: grpPrivGrpMember
//         Ties group privileges to groups for whose members the privilege
//         is granted.
//
func (s *Storage) newGrpPrivGrpMember() *grpPrivGrpMember {
	o := grpPrivGrpMember{}
	return &o
}
