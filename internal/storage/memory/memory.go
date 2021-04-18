// brigolagecms/pkg/storage/memory/memory.go
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

// Package memory implements the in-memory data store for BrigolageCMS.
package memory

type Storage struct {
	sequences struct {
		action                   *sequence
		actionType               *sequence
		addr                     *sequence
		addrPart                 *sequence
		addrPartType             *sequence
		alert                    *sequence
		alertType                *sequence
		alertTypeMember          *sequence
		alertTypeRule            *sequence
		alerted                  *sequence
		attrAction               *sequence
		attrActionMeta           *sequence
		attrActionVal            *sequence
		attrCategory             *sequence
		attrCategoryMeta         *sequence
		attrCategoryVal          *sequence
		attrElementType          *sequence
		attrElementTypeMeta      *sequence
		attrElementTypeVal       *sequence
		attrFieldType            *sequence
		attrFieldTypeMeta        *sequence
		attrFieldTypeVal         *sequence
		attrGrp                  *sequence
		attrGrpMeta              *sequence
		attrGrpVal               *sequence
		attrMember               *sequence
		attrMemberMeta           *sequence
		attrMemberVal            *sequence
		category                 *sequence
		categoryMember           *sequence
		class                    *sequence
		contact                  *sequence
		contactValue             *sequence
		contribTypeMember        *sequence
		desk                     *sequence
		deskMember               *sequence
		destMember               *sequence
		elementLanguage          *sequence
		elementType              *sequence
		elementTypeMember        *sequence
		elementTypeOutputChannel *sequence
		elementTypeSite          *sequence
		event                    *sequence
		eventAttr                *sequence
		eventType                *sequence
		eventTypeAttr            *sequence
		fieldType                *sequence
		grp                      *sequence
		grpMember                *sequence
		job                      *sequence
		jobMember                *sequence
		keyword                  *sequence
		keywordMember            *sequence
		mediaContributor         *sequence
		mediaElement             *sequence
		mediaField               *sequence
		mediaFields              *sequence
		mediaMember              *sequence
		mediaType                *sequence
		mediaTypeExt             *sequence
		mediaTypeMember          *sequence
		mediaURI                 *sequence
		member                   *sequence
		org                      *sequence
		orgMember                *sequence
		outputChannel            *sequence
		outputChannelInclude     *sequence
		outputChannelMember      *sequence
		person                   *sequence
		personMember             *sequence
		personOrg                *sequence
		pref                     *sequence
		prefMember               *sequence
		priv                     *sequence
		resource                 *sequence
		server                   *sequence
		serverType               *sequence
		siteMember               *sequence
		source                   *sequence
		sourceMember             *sequence
		story                    *sequence
		storyCategory            *sequence
		storyContributor         *sequence
		storyElement             *sequence
		storyField               *sequence
		storyInstance            *sequence
		storyMember              *sequence
		storyURI                 *sequence
		subelementType           *sequence
		template                 *sequence
		templateInstance         *sequence
		templateMember           *sequence
		userMember               *sequence
		usrPref                  *sequence
		workflow                 *sequence
		workflowMember           *sequence
	}
	tables struct {
		action                   map[int]*action
		actionType               map[int]*actionType
		actionTypeMediaType      map[int]*actionTypeMediaType
		addr                     map[int]*addr
		addrPart                 map[int]*addrPart
		addrPartType             map[int]*addrPartType
		alert                    map[int]*alert
		alertType                map[int]*alertType
		alertTypeGrpContact      map[int]*alertTypeGrpContact
		alertTypeUsrContact      map[int]*alertTypeUsrContact
		alertTypeMember          map[int]*alertTypeMember
		alertTypeRule            map[int]*alertTypeRule
		alerted                  map[int]*alerted
		alertedContactValue      map[int]*alertedContactValue
		atType                   *atTypeTable
		atTypeMember             *atTypeMemberTable
		attrAction               map[int]*attrAction
		attrActionMeta           map[int]*attrActionMeta
		attrActionVal            map[int]*attrActionVal
		attrCategory             map[int]*attrCategory
		attrCategoryMeta         map[int]*attrCategoryMeta
		attrCategoryVal          map[int]*attrCategoryVal
		attrElementType          map[int]*attrElementType
		attrElementTypeMeta      map[int]*attrElementTypeMeta
		attrElementTypeVal       map[int]*attrElementTypeVal
		attrFieldType            map[int]*attrFieldType
		attrFieldTypeMeta        map[int]*attrFieldTypeMeta
		attrFieldTypeVal         map[int]*attrFieldTypeVal
		attrGrp                  map[int]*attrGrp
		attrGrpMeta              map[int]*attrGrpMeta
		attrGrpVal               map[int]*attrGrpVal
		attrMember               map[int]*attrMember
		attrMemberMeta           map[int]*attrMemberMeta
		attrMemberVal            map[int]*attrMemberVal
		category                 map[int]*category
		categoryKeyword          map[int]*categoryKeyword
		categoryMember           map[int]*categoryMember
		class                    map[int]*class
		contact                  map[int]*contact
		contactValue             map[int]*contactValue
		contribTypeMember        map[int]*contribTypeMember
		desk                     map[int]*desk
		deskMember               map[int]*deskMember
		destMember               map[int]*destMember
		elementType              map[int]*elementType
		elementTypeOutputChannel map[int]*elementTypeOutputChannel
		elementTypeSite          map[int]*elementTypeSite
		elementTypeMember        map[int]*elementTypeMember
		event                    map[int]*event
		eventAttr                map[int]*eventAttr
		eventType                map[int]*eventType
		eventTypeAttr            map[int]*eventTypeAttr
		fieldType                map[int]*fieldType
		grp                      map[int]*grp
		grpMember                map[int]*grpMember
		grpPriv                  map[int]*grpPriv
		grpPrivGrpMember         map[int]*grpPrivGrpMember
		job                      map[int]*job
		jobResource              map[int]*jobResource
		jobServerType            map[int]*jobServerType
		jobMember                map[int]*jobMember
		keyword                  map[int]*keyword
		keywordMember            map[int]*keywordMember
		media                    *mediaTable
		mediaContributor         map[int]*mediaContributor
		mediaElement             map[int]*mediaElement
		mediaField               map[int]*mediaField
		mediaFields              map[int]*mediaFields
		mediaInstance            *mediaInstanceTable
		mediaKeyword             map[int]*mediaKeyword
		mediaMember              map[int]*mediaMember
		mediaOutputChannel       map[int]*mediaOutputChannel
		mediaResource            map[int]*mediaResource
		mediaType                map[int]*mediaType
		mediaTypeExt             map[int]*mediaTypeExt
		mediaTypeMember          map[int]*mediaTypeMember
		mediaUri                 map[int]*mediaURI
		member                   map[int]*member
		org                      map[int]*org
		orgMember                map[int]*orgMember
		outputChannel            map[int]*outputChannel
		outputChannelInclude     map[int]*outputChannelInclude
		outputChannelMember      map[int]*outputChannelMember
		person                   map[int]*person
		personContactValue       map[int]*personContactValue
		personMember             map[int]*personMember
		personOrg                map[int]*personOrg
		personOrgAddr            map[int]*personOrgAddr
		pref                     map[int]*pref
		prefMember               map[int]*prefMember
		prefOpt                  map[int]*prefOpt
		resource                 map[int]*resource
		server                   map[int]*server
		serverType               map[int]*serverType
		serverTypeOutputChannel  map[int]*serverTypeOutputChannel
		site                     map[int]*site
		siteMember               map[int]*siteMember
		source                   map[int]*source
		sourceMember             map[int]*sourceMember
		story                    map[int]*story
		storyCategory            map[int]*storyCategory
		storyContributor         map[int]*storyContributor
		storyOutputChannel       map[int]*storyOutputChannel
		storyResource            map[int]*storyResource
		storyElement             map[int]*storyElement
		storyField               map[int]*storyField
		storyInstance            map[int]*storyInstance
		storyKeyword             map[int]*storyKeyword
		storyMember              map[int]*storyMember
		storyUri                 map[int]*storyURI
		subelementType           map[int]*subelementType
		template                 map[int]*template
		templateInstance         map[int]*templateInstance
		templateMember           map[int]*templateMember
		userMember               map[int]*userMember
		usr                      map[int]*usr
		usrPref                  map[int]*usrPref
		workflow                 map[int]*workflow
		workflowMember           map[int]*workflowMember
	}
}

func New() (*Storage, error) {
	var s Storage

	s.sequences.action = newSequence(1024)
	s.sequences.actionType = newSequence(1024)
	s.sequences.addr = newSequence(1024)
	s.sequences.addrPart = newSequence(1024)
	s.sequences.addrPartType = newSequence(1024)
	s.sequences.alert = newSequence(1024)
	s.sequences.alertType = newSequence(1024)
	s.sequences.alertTypeMember = newSequence(1024)
	s.sequences.alertTypeRule = newSequence(1024)
	s.sequences.alerted = newSequence(1024)
	s.sequences.attrAction = newSequence(1024)
	s.sequences.attrActionMeta = newSequence(1024)
	s.sequences.attrActionVal = newSequence(1024)
	s.sequences.attrCategory = newSequence(1024)
	s.sequences.attrCategoryMeta = newSequence(1024)
	s.sequences.attrCategoryVal = newSequence(1024)
	s.sequences.attrElementType = newSequence(1024)
	s.sequences.attrElementTypeMeta = newSequence(1024)
	s.sequences.attrElementTypeVal = newSequence(1024)
	s.sequences.attrFieldType = newSequence(1024)
	s.sequences.attrFieldTypeMeta = newSequence(1024)
	s.sequences.attrFieldTypeVal = newSequence(1024)
	s.sequences.attrGrp = newSequence(1024)
	s.sequences.attrGrpMeta = newSequence(1024)
	s.sequences.attrGrpVal = newSequence(1024)
	s.sequences.attrMember = newSequence(1024)
	s.sequences.attrMemberMeta = newSequence(1024)
	s.sequences.attrMemberVal = newSequence(1024)
	s.sequences.category = newSequence(1024)
	s.sequences.categoryMember = newSequence(1024)
	s.sequences.class = newSequence(1024)
	s.sequences.contact = newSequence(1024)
	s.sequences.contactValue = newSequence(1024)
	s.sequences.contribTypeMember = newSequence(1024)
	s.sequences.desk = newSequence(1024)
	s.sequences.deskMember = newSequence(1024)
	s.sequences.destMember = newSequence(1024)
	s.sequences.elementLanguage = newSequence(1024)
	s.sequences.elementType = newSequence(1024)
	s.sequences.elementTypeMember = newSequence(1024)
	s.sequences.elementTypeOutputChannel = newSequence(1024)
	s.sequences.elementTypeSite = newSequence(1024)
	s.sequences.event = newSequence(1024)
	s.sequences.eventAttr = newSequence(1024)
	s.sequences.eventType = newSequence(1024)
	s.sequences.eventTypeAttr = newSequence(1024)
	s.sequences.fieldType = newSequence(1024)
	s.sequences.grp = newSequence(1024)
	s.sequences.grpMember = newSequence(1024)
	s.sequences.job = newSequence(1024)
	s.sequences.jobMember = newSequence(1024)
	s.sequences.keyword = newSequence(1024)
	s.sequences.keywordMember = newSequence(1024)
	s.sequences.mediaContributor = newSequence(1024)
	s.sequences.mediaElement = newSequence(1024)
	s.sequences.mediaField = newSequence(1024)
	s.sequences.mediaFields = newSequence(1024)
	s.sequences.mediaMember = newSequence(1024)
	s.sequences.mediaType = newSequence(1024)
	s.sequences.mediaTypeExt = newSequence(1024)
	s.sequences.mediaTypeMember = newSequence(1024)
	s.sequences.mediaURI = newSequence(1024)
	s.sequences.member = newSequence(1024)
	s.sequences.org = newSequence(1024)
	s.sequences.orgMember = newSequence(1024)
	s.sequences.outputChannel = newSequence(1024)
	s.sequences.outputChannelInclude = newSequence(1024)
	s.sequences.outputChannelMember = newSequence(1024)
	s.sequences.person = newSequence(1024)
	s.sequences.personMember = newSequence(1024)
	s.sequences.personOrg = newSequence(1024)
	s.sequences.pref = newSequence(1024)
	s.sequences.prefMember = newSequence(1024)
	s.sequences.priv = newSequence(1024)
	s.sequences.resource = newSequence(1024)
	s.sequences.server = newSequence(1024)
	s.sequences.serverType = newSequence(1024)
	s.sequences.siteMember = newSequence(1024)
	s.sequences.source = newSequence(1024)
	s.sequences.sourceMember = newSequence(1024)
	s.sequences.story = newSequence(1024)
	s.sequences.storyCategory = newSequence(1024)
	s.sequences.storyContributor = newSequence(1024)
	s.sequences.storyElement = newSequence(1024)
	s.sequences.storyField = newSequence(1024)
	s.sequences.storyInstance = newSequence(1024)
	s.sequences.storyMember = newSequence(1024)
	s.sequences.storyURI = newSequence(1024)
	s.sequences.subelementType = newSequence(1024)
	s.sequences.template = newSequence(1024)
	s.sequences.templateInstance = newSequence(1024)
	s.sequences.templateMember = newSequence(1024)
	s.sequences.userMember = newSequence(1024)
	s.sequences.usrPref = newSequence(1024)
	s.sequences.workflow = newSequence(1024)
	s.sequences.workflowMember = newSequence(1024)

	s.tables.action = make(map[int]*action)
	s.tables.actionType = make(map[int]*actionType)
	s.tables.actionTypeMediaType = make(map[int]*actionTypeMediaType)
	s.tables.addr = make(map[int]*addr)
	s.tables.addrPart = make(map[int]*addrPart)
	s.tables.addrPartType = make(map[int]*addrPartType)
	s.tables.alert = make(map[int]*alert)
	s.tables.alertType = make(map[int]*alertType)
	s.tables.alertTypeGrpContact = make(map[int]*alertTypeGrpContact)
	s.tables.alertTypeUsrContact = make(map[int]*alertTypeUsrContact)
	s.tables.alertTypeMember = make(map[int]*alertTypeMember)
	s.tables.alertTypeRule = make(map[int]*alertTypeRule)
	s.tables.alerted = make(map[int]*alerted)
	s.tables.alertedContactValue = make(map[int]*alertedContactValue)
	s.tables.atType = &atTypeTable{seq: newSequence(1024), rows: make(map[int]*atType)}
	s.tables.atTypeMember = &atTypeMemberTable{seq: newSequence(1024), rows: make(map[int]*atTypeMember)}
	s.tables.attrAction = make(map[int]*attrAction)
	s.tables.attrActionMeta = make(map[int]*attrActionMeta)
	s.tables.attrActionVal = make(map[int]*attrActionVal)
	s.tables.attrCategory = make(map[int]*attrCategory)
	s.tables.attrCategoryMeta = make(map[int]*attrCategoryMeta)
	s.tables.attrCategoryVal = make(map[int]*attrCategoryVal)
	s.tables.attrElementType = make(map[int]*attrElementType)
	s.tables.attrElementTypeMeta = make(map[int]*attrElementTypeMeta)
	s.tables.attrElementTypeVal = make(map[int]*attrElementTypeVal)
	s.tables.attrFieldType = make(map[int]*attrFieldType)
	s.tables.attrFieldTypeMeta = make(map[int]*attrFieldTypeMeta)
	s.tables.attrFieldTypeVal = make(map[int]*attrFieldTypeVal)
	s.tables.attrGrp = make(map[int]*attrGrp)
	s.tables.attrGrpMeta = make(map[int]*attrGrpMeta)
	s.tables.attrGrpVal = make(map[int]*attrGrpVal)
	s.tables.attrMember = make(map[int]*attrMember)
	s.tables.attrMemberMeta = make(map[int]*attrMemberMeta)
	s.tables.attrMemberVal = make(map[int]*attrMemberVal)
	s.tables.category = make(map[int]*category)
	s.tables.categoryKeyword = make(map[int]*categoryKeyword)
	s.tables.categoryMember = make(map[int]*categoryMember)
	s.tables.class = make(map[int]*class)
	s.tables.contact = make(map[int]*contact)
	s.tables.contactValue = make(map[int]*contactValue)
	s.tables.contribTypeMember = make(map[int]*contribTypeMember)
	s.tables.desk = make(map[int]*desk)
	s.tables.deskMember = make(map[int]*deskMember)
	s.tables.destMember = make(map[int]*destMember)
	s.tables.elementType = make(map[int]*elementType)
	s.tables.elementTypeOutputChannel = make(map[int]*elementTypeOutputChannel)
	s.tables.elementTypeSite = make(map[int]*elementTypeSite)
	s.tables.elementTypeMember = make(map[int]*elementTypeMember)
	s.tables.event = make(map[int]*event)
	s.tables.eventAttr = make(map[int]*eventAttr)
	s.tables.eventType = make(map[int]*eventType)
	s.tables.eventTypeAttr = make(map[int]*eventTypeAttr)
	s.tables.fieldType = make(map[int]*fieldType)
	s.tables.grp = make(map[int]*grp)
	s.tables.grpMember = make(map[int]*grpMember)
	s.tables.grpPriv = make(map[int]*grpPriv)
	s.tables.grpPrivGrpMember = make(map[int]*grpPrivGrpMember)
	s.tables.job = make(map[int]*job)
	s.tables.jobResource = make(map[int]*jobResource)
	s.tables.jobServerType = make(map[int]*jobServerType)
	s.tables.jobMember = make(map[int]*jobMember)
	s.tables.keyword = make(map[int]*keyword)
	s.tables.keywordMember = make(map[int]*keywordMember)
	s.tables.media = &mediaTable{seq: newSequence(1024), rows: make(map[int]*media)}
	s.tables.mediaContributor = make(map[int]*mediaContributor)
	s.tables.mediaOutputChannel = make(map[int]*mediaOutputChannel)
	s.tables.mediaResource = make(map[int]*mediaResource)
	s.tables.mediaElement = make(map[int]*mediaElement)
	s.tables.mediaField = make(map[int]*mediaField)
	s.tables.mediaFields = make(map[int]*mediaFields)
	s.tables.mediaInstance = &mediaInstanceTable{seq: newSequence(1024), rows: make(map[int]*mediaInstance)}
	s.tables.mediaKeyword = make(map[int]*mediaKeyword)
	s.tables.mediaMember = make(map[int]*mediaMember)
	s.tables.mediaType = make(map[int]*mediaType)
	s.tables.mediaTypeExt = make(map[int]*mediaTypeExt)
	s.tables.mediaTypeMember = make(map[int]*mediaTypeMember)
	s.tables.mediaUri = make(map[int]*mediaURI)
	s.tables.member = make(map[int]*member)
	s.tables.org = make(map[int]*org)
	s.tables.orgMember = make(map[int]*orgMember)
	s.tables.outputChannel = make(map[int]*outputChannel)
	s.tables.outputChannelInclude = make(map[int]*outputChannelInclude)
	s.tables.outputChannelMember = make(map[int]*outputChannelMember)
	s.tables.person = make(map[int]*person)
	s.tables.personContactValue = make(map[int]*personContactValue)
	s.tables.personMember = make(map[int]*personMember)
	s.tables.personOrg = make(map[int]*personOrg)
	s.tables.personOrgAddr = make(map[int]*personOrgAddr)
	s.tables.pref = make(map[int]*pref)
	s.tables.prefMember = make(map[int]*prefMember)
	s.tables.prefOpt = make(map[int]*prefOpt)
	s.tables.resource = make(map[int]*resource)
	s.tables.server = make(map[int]*server)
	s.tables.serverType = make(map[int]*serverType)
	s.tables.serverTypeOutputChannel = make(map[int]*serverTypeOutputChannel)
	s.tables.site = make(map[int]*site)
	s.tables.siteMember = make(map[int]*siteMember)
	s.tables.source = make(map[int]*source)
	s.tables.sourceMember = make(map[int]*sourceMember)
	s.tables.story = make(map[int]*story)
	s.tables.storyCategory = make(map[int]*storyCategory)
	s.tables.storyContributor = make(map[int]*storyContributor)
	s.tables.storyOutputChannel = make(map[int]*storyOutputChannel)
	s.tables.storyResource = make(map[int]*storyResource)
	s.tables.storyElement = make(map[int]*storyElement)
	s.tables.storyField = make(map[int]*storyField)
	s.tables.storyInstance = make(map[int]*storyInstance)
	s.tables.storyKeyword = make(map[int]*storyKeyword)
	s.tables.storyMember = make(map[int]*storyMember)
	s.tables.storyUri = make(map[int]*storyURI)
	s.tables.subelementType = make(map[int]*subelementType)
	s.tables.template = make(map[int]*template)
	s.tables.templateInstance = make(map[int]*templateInstance)
	s.tables.templateMember = make(map[int]*templateMember)
	s.tables.userMember = make(map[int]*userMember)
	s.tables.usr = make(map[int]*usr)
	s.tables.usrPref = make(map[int]*usrPref)
	s.tables.workflow = make(map[int]*workflow)
	s.tables.workflowMember = make(map[int]*workflowMember)

	return &s, nil
}
