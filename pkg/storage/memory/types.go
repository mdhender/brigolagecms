// brigolagecms/pkg/storage/memory/types.go
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

// This file contains all of the POSTGRES tables and sequences.
// None of the functions, indexes, or other objects were extracted.
//
// The objects were extracted using `ddltogo`. Table, Sequences,
// and Column names were processed through `vars.sed` to match
// Go's linter's suggestions. The comments are left intact for
// personal convenience.

import "time"

//  Project: Bricolage
//
//  Author: David Wheeler <david@justatheory.com>
//
//  This DDL is for the creation of universal stuff needed by other DDLs, such as
//  functions.
//
//  Functions.
//
//  This function allows us to create UNIQUE indices that combine a lowercased
//  TEXT (or VARCHAR) column with an INTEGER column. See Bric/Util/AlertType.sql
//  for an example.
//   21: /* XXX Once we require 7.4 or later (2.0?), dump the function and aggregate
//   21:  * below in favor of this aggregate:
//   21:
//   21: CREATE AGGREGATE array_accum (
//   21:     sfunc    = array_append,
//   21:     basetype = anyelement,
//   21:     stype    = anyarray,
//   21:     initcond = '{}'
//   21: );
//   21:
//   21:  * Then change the usage from group_concat(foo) to
//   21:  * array_to_string(array_accum(foo), ' ')
//   21:
//   21: */
//  This function is used to append a space followed by a number to a TEXT
//  string. It is used primarily for the group_concat aggregate (below). We omit
//  the ID 0 because it is a hidden, secret group to which permissions do not
//  apply.
//  This aggregate is designed to concatenate all of the IDs that would
//  otherwise cause a query to return multiple rows into a single value
//  with each ID separated by a space. This makes it easy for us to pull
//  out the list of IDs using split, _and_ keeps the entire contents of
//  an object in a single row, thus also enabling the use of OFFSET and
//  LIMIT.
//   63: /*
//   63: -- This is a temporary compatibility measure.
//   63: CREATE FUNCTION int_to_boolean(integer) RETURNS boolean
//   63:   AS 'select case when $1 = 0 then false else true end'
//   63: LANGUAGE 'sql' IMMUTABLE;
//   63:
//   63: CREATE CAST (integer AS boolean)
//   63:   WITH FUNCTION int_to_boolean(integer) AS IMPLICIT;
//   63: */
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  This is the SQL that will create the element table.
//  It is related to the Bric::ElementType class.
//
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the element table
// create sequence seqAtType
// start 1024
type seqAtType struct {
	Val int
}

func newSeqAtType() *seqAtType {
	return &seqAtType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAtType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAtType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAtTypeMember
// start 1024
type seqAtTypeMember struct {
	Val int
}

func newSeqAtTypeMember() *seqAtTypeMember {
	return &seqAtTypeMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAtTypeMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAtTypeMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: element
//
//  Description:    The table that holds the information for a given asset type.
//          Holds name and description information and is references by
//         elementType rows.
//
// create table atType
// constraint pk_atTypeID
//    pk(id)
type atType struct {
	id           int    // integer
	name         string // varchar(64)
	description  string // varchar(256)
	topLevel     bool
	paginated    bool
	fixedURL     bool
	relatedStory bool
	relatedMedia bool
	media        bool
	bizClassID   int // integer
	active       bool
}

func newAtType() *atType {
	o := atType{
		id:           0, // seqAtType.NextVal(),
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

//
//  TABLE: elementTypeMember
//
// create table atTypeMember
// constraint pk_atTypeMember__id
//    pk(id)
type atTypeMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newAtTypeMember() *atTypeMember {
	o := atTypeMember{
		id: 0, // seqAtTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  The sql representation of Media Assets.
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the media table
// create sequence seqMedia
// start 1024
type seqMedia struct {
	Val int
}

func newSeqMedia() *seqMedia {
	return &seqMedia{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMedia) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMedia) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqMediaInstance
// start 1024
type seqMediaInstance struct {
	Val int
}

func newSeqMediaInstance() *seqMediaInstance {
	return &seqMediaInstance{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaInstance) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaInstance) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique ids for the mediaContributor table
// create sequence seqMediaContributor
// start 1024
type seqMediaContributor struct {
	Val int
}

func newSeqMediaContributor() *seqMediaContributor {
	return &seqMediaContributor{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaContributor) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaContributor) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the mediaMember table
// create sequence seqMediaMember
// start 1024
type seqMediaMember struct {
	Val int
}

func newSeqMediaMember() *seqMediaMember {
	return &seqMediaMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqMediaFields
// start 1024
type seqMediaFields struct {
	Val int
}

func newSeqMediaFields() *seqMediaFields {
	return &seqMediaFields{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaFields) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaFields) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the mediaURI table
// create sequence seqMediaURI
// start 1024
type seqMediaURI struct {
	Val int
}

func newSeqMediaURI() *seqMediaURI {
	return &seqMediaURI{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaURI) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaURI) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table media
//
//  Description: The Media table this houses the data for a given media asset
//                               and its related asset_version_data
//
// create table media
// constraint pk_mediaID
//    pk(id)
type media struct {
	id               int       // integer
	uuid             []byte    // text, like a blob
	elementTypeID    int       // integer
	currentVersion   int       // integer
	publishedVersion int       // integer
	usrID            int       // integer
	firstPublishDate time.Time // timestamp
	publishDate      time.Time // timestamp
	workflowID       int       // integer
	deskID           int       // integer
	publishStatus    bool
	active           bool
	siteID           int // integer
	aliasID          int // integer
}

func newMedia() *media {
	o := media{
		id:            0, // seqMedia.NextVal(),
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
//
// create table mediaInstance
// constraint pk_mediaInstanceID
//    pk(id)
type mediaInstance struct {
	id          int       // integer
	name        string    // varchar(256)
	description string    // varchar(1024)
	mediaID     int       // integer
	sourceID    int       // integer
	priority    int       // integer
	usrID       int       // integer
	version     int       // integer
	expireDate  time.Time // timestamp
	categoryID  int       // integer
	mediaTypeID int       // integer
	primaryOcID int       // integer
	fileSize    int       // integer
	fileName    string    // varchar(256)
	location    string    // varchar(256)
	uri         string    // varchar(256)
	coverDate   time.Time // timestamp
	note        []byte    // text, like a blob
	checkedOut  bool
}

func newMediaInstance() *mediaInstance {
	o := mediaInstance{
		id:         0, // seqMediaInstance.NextVal(),
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
// create table mediaURI
// constraint pk_mediaURI__id
//    pk(id)
type mediaURI struct {
	id      int    // integer
	mediaID int    // integer
	siteID  int    // integer
	uri     []byte // text, like a blob
}

func newMediaURI() *mediaURI {
	o := mediaURI{
		id: 0, // seqMediaURI.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaOutputChannel
//
//  Description: Mapping Table between stories and output channels.
//
//
// create table mediaOutputChannel
// constraint pk_media_outputChannel
//    pk(mediaInstanceID)
//    pk(outputChannelID)
type mediaOutputChannel struct {
	mediaInstanceID int // integer
	outputChannelID int // integer
}

func newMediaOutputChannel() *mediaOutputChannel {
	o := mediaOutputChannel{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaFields
//
//  Description: A mapping table between Media classes and functions that
//                               Will be run against uploaded files
//
// create table mediaFields
// constraint pk_mediaFields__id
//    pk(id)
type mediaFields struct {
	id           int    // integer
	bizPkg       int    // integer
	name         string // varchar(32)
	functionName string // varchar(256)
	active       bool
}

func newMediaFields() *mediaFields {
	o := mediaFields{
		id:     0, // seqMediaFields.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table mediaContributor
//
//  Description: mapping tables between media instances and contributors
//
//
// create table mediaContributor
// constraint pk_media_categoryID
//    pk(id)
type mediaContributor struct {
	id              int    // integer
	mediaInstanceID int    // integer
	memberID        int    // integer
	place           int    // integer
	role            string // varchar(256)
}

func newMediaContributor() *mediaContributor {
	o := mediaContributor{
		id: 0, // seqMediaContributor.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaMember
//
//  Description: The link between media objects and member objects
//
// create table mediaMember
// constraint pk_mediaMember__id
//    pk(id)
type mediaMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newMediaMember() *mediaMember {
	o := mediaMember{
		id: 0, // seqMediaMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  media
//  mediaInstance
//  mediaURI
//  mediaOutputChannel
// mediaContributor
//  mediaMember.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//
//  This is the SQL representation of the story object
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the story table
// create sequence seqStory
// start 1024
type seqStory struct {
	Val int
}

func newSeqStory() *seqStory {
	return &seqStory{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStory) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStory) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the storyInstance table
// create sequence seqStoryInstance
// start 1024
type seqStoryInstance struct {
	Val int
}

func newSeqStoryInstance() *seqStoryInstance {
	return &seqStoryInstance{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryInstance) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryInstance) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the storyCategory mapping table
// create sequence seqStoryCategory
// start 1024
type seqStoryCategory struct {
	Val int
}

func newSeqStoryCategory() *seqStoryCategory {
	return &seqStoryCategory{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryCategory) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryCategory) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique ids for the storyContributor table
// create sequence seqStoryContributor
// start 1024
type seqStoryContributor struct {
	Val int
}

func newSeqStoryContributor() *seqStoryContributor {
	return &seqStoryContributor{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryContributor) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryContributor) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the storyURI table
// create sequence seqStoryURI
// start 1024
type seqStoryURI struct {
	Val int
}

func newSeqStoryURI() *seqStoryURI {
	return &seqStoryURI{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryURI) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryURI) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: story
//
//  Description: The story properties. Versioning info might get added here and
//               the rights info might get removed. It is also possible that
//               the asset type field will need a cascading delete.
// create table story
// constraint pk_storyID
//    pk(id)
type story struct {
	id               int       // integer
	uuid             []byte    // text, like a blob
	usrID            int       // integer
	elementTypeID    int       // integer
	firstPublishDate time.Time // timestamp
	publishDate      time.Time // timestamp
	currentVersion   int       // integer
	publishedVersion int       // integer
	workflowID       int       // integer
	deskID           int       // integer
	publishStatus    bool
	active           bool
	siteID           int // integer
	aliasID          int // integer
}

func newStory() *story {
	o := story{
		id:            0, // seqStory.NextVal(),
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
// create table storyInstance
// constraint pk_storyInstanceID
//    pk(id)
type storyInstance struct {
	id          int       // integer
	name        string    // varchar(256)
	description string    // varchar(1024)
	storyID     int       // integer
	sourceID    int       // integer
	primaryURI  string    // varchar(128)
	priority    int       // integer
	version     int       // integer
	expireDate  time.Time // timestamp
	usrID       int       // integer
	slug        string    // varchar(64)
	primaryOcID int       // integer
	coverDate   time.Time // timestamp
	note        []byte    // text, like a blob
	checkedOut  bool
}

func newStoryInstance() *storyInstance {
	o := storyInstance{
		id:         0, // seqStoryInstance.NextVal(),
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
// create table storyURI
// constraint pk_storyURI__id
//    pk(id)
type storyURI struct {
	id      int    // integer
	storyID int    // integer
	siteID  int    // integer
	uri     []byte // text, like a blob
}

func newStoryURI() *storyURI {
	o := storyURI{
		id: 0, // seqStoryURI.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyOutputChannel
//
//  Description: Mapping Table between stories and output channels.
//
//
// create table storyOutputChannel
// constraint pk_story_outputChannel
//    pk(storyInstanceID)
//    pk(outputChannelID)
type storyOutputChannel struct {
	storyInstanceID int // integer
	outputChannelID int // integer
}

func newStoryOutputChannel() *storyOutputChannel {
	o := storyOutputChannel{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyCategory
//
//  Description: Mapping Table between Stories and categories
//
//
// create table storyCategory
// constraint pk_story_categoryID
//    pk(id)
type storyCategory struct {
	id              int // integer
	storyInstanceID int // integer
	categoryID      int // integer
	main            bool
}

func newStoryCategory() *storyCategory {
	o := storyCategory{
		id:   0, // seqStoryCategory.NextVal(),
		main: false,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table storyContributor
//
//  Description: mapping tables between story instances and contributors
//
//
// create table storyContributor
// constraint pk_story_categoryID
//    pk(id)
type storyContributor struct {
	id              int    // integer
	storyInstanceID int    // integer
	memberID        int    // integer
	place           int    // integer
	role            string // varchar(256)
}

func newStoryContributor() *storyContributor {
	o := storyContributor{
		id: 0, // seqStoryContributor.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  story
//  storyInstance
//  storyURI
//  storyCategory
//  storyOutputChannel
// storyContributor
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  The sql that will hold all the template asset information.
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique ids for the template table
// create sequence seqTemplate
// start 1024
type seqTemplate struct {
	Val int
}

func newSeqTemplate() *seqTemplate {
	return &seqTemplate{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqTemplate) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqTemplate) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqTemplateInstance
// start 1024
type seqTemplateInstance struct {
	Val int
}

func newSeqTemplateInstance() *seqTemplateInstance {
	return &seqTemplateInstance{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqTemplateInstance) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqTemplateInstance) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the storyMember table
// create sequence seqTemplateMember
// start 1024
type seqTemplateMember struct {
	Val int
}

func newSeqTemplateMember() *seqTemplateMember {
	return &seqTemplateMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqTemplateMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqTemplateMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table template
//
//  Description: The table that holds all the template info
//
//
// create table template
// constraint pk_templateID
//    pk(id)
type template struct {
	id               int    // integer
	usrID            int    // integer
	outputChannelID  int    // integer
	tplateType       int    // integer
	elementTypeID    int    // integer
	fileName         []byte // text, like a blob
	currentVersion   int    // integer
	workflowID       int    // integer
	deskID           int    // integer
	publishedVersion int    // integer
	deployStatus     bool
	deployDate       time.Time // timestamp
	active           bool
	siteID           int // integer
}

func newTemplate() *template {
	o := template{
		id:           0, // seqTemplate.NextVal(),
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
// create table templateInstance
// constraint pk_templateInstance__id
//    pk(id)
type templateInstance struct {
	id          int       // integer
	name        string    // varchar(256)
	description string    // varchar(1024)
	priority    int       // integer
	templateID  int       // integer
	categoryID  int       // integer
	version     int       // integer
	usrID       int       // integer
	fileName    []byte    // text, like a blob
	data        []byte    // text, like a blob
	expireDate  time.Time // timestamp
	note        []byte    // text, like a blob
	checkedOut  bool
}

func newTemplateInstance() *templateInstance {
	o := templateInstance{
		id:         0, // seqTemplateInstance.NextVal(),
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
// create table templateMember
// constraint pk_templateMember__id
//    pk(id)
type templateMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newTemplateMember() *templateMember {
	o := templateMember{
		id: 0, // seqTemplateMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  template.
//  templateInstance.
//  templateMember.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  This SQL creates the tables necessary for the attribute object.  This file
//  applies to attributes on the Bric::Person class.  Any other classes that
//  require attributes need only duplicate these tables, changing 'person' to
//  the correct class name.  Class names may be shortened to ensure that the
//  resulting table names are under the oracle 30 character name limit as long
//  as the resulting shortened class name is unique.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the category table
// create sequence seqCategory
// start 1024
type seqCategory struct {
	Val int
}

func newSeqCategory() *seqCategory {
	return &seqCategory{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqCategory) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqCategory) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the categoryMember table
// create sequence seqCategoryMember
// start 1024
type seqCategoryMember struct {
	Val int
}

func newSeqCategoryMember() *seqCategoryMember {
	return &seqCategoryMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqCategoryMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqCategoryMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the attrCategory table
// create sequence seqAttrCategory
// start 1024
type seqAttrCategory struct {
	Val int
}

func newSeqAttrCategory() *seqAttrCategory {
	return &seqAttrCategory{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrCategory) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrCategory) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for each attrCategoryVal table
// create sequence seqAttrCategoryVal
// start 1024
type seqAttrCategoryVal struct {
	Val int
}

func newSeqAttrCategoryVal() *seqAttrCategoryVal {
	return &seqAttrCategoryVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrCategoryVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrCategoryVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the category_meta table
// create sequence seqAttrCategoryMeta
// start 1024
type seqAttrCategoryMeta struct {
	Val int
}

func newSeqAttrCategoryMeta() *seqAttrCategoryMeta {
	return &seqAttrCategoryMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrCategoryMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrCategoryMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: category
//
//  Description: The category table
//
// create table category
// constraint pk_categoryID
//    pk(id)
type category struct {
	id          int    // integer
	siteID      int    // integer
	directory   string // varchar(128)
	uri         string // varchar(256)
	name        string // varchar(128)
	description string // varchar(256)
	parentID    int    // integer
	assetGrpID  int    // integer
	active      bool
}

func newCategory() *category {
	o := category{
		id:     0, // seqCategory.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: categoryMember
//
//  Description: The link between desk objects and member objects
//
// create table categoryMember
// constraint pk_categoryMember__id
//    pk(id)
type categoryMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newCategoryMember() *categoryMember {
	o := categoryMember{
		id: 0, // seqCategoryMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategory
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its category ID and an attribute name.
// create table attrCategory
// constraint pk_attrCategory__id
//    pk(id)
type attrCategory struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrCategory() *attrCategory {
	o := attrCategory{
		id:     0, // seqAttrCategory.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategoryVal
//
//  Description: A table to hold attribute values.
// create table attrCategoryVal
// constraint pk_attrCategoryVal__id
//    pk(id)
type attrCategoryVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrCategoryVal() *attrCategoryVal {
	o := attrCategoryVal{
		id:     0, // seqAttrCategoryVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrCategoryMeta
//
//  Description: A table to represent metadata on types of attributes.
// create table attrCategoryMeta
// constraint pk_attrCategoryMeta__id
//    pk(id)
type attrCategoryMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  string // varchar(2048)
	active bool
}

func newAttrCategoryMeta() *attrCategoryMeta {
	o := attrCategoryMeta{
		id:     0, // seqAttrCategoryMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//  This DDL creates the basic tables for all Bric::BC::Contact objects.
//
//  SEQUENCES.
//
// create sequence seqContact
// start 1024
type seqContact struct {
	Val int
}

func newSeqContact() *seqContact {
	return &seqContact{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqContact) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqContact) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqContactValue
// start 1024
type seqContactValue struct {
	Val int
}

func newSeqContactValue() *seqContactValue {
	return &seqContactValue{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqContactValue) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqContactValue) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: contact
//
// create table contact
// constraint pk_contactID
//    pk(id)
type contact struct {
	id          int    // integer
	ttypee      string // varchar(64)
	description string // varchar(256)
	active      bool
	alertable   bool
}

func newContact() *contact {
	o := contact{
		id:        0, // seqContact.NextVal(),
		active:    true,
		alertable: false,
	}
	return &o
}

//
//  TABLE: contactValue
//
// create table contactValue
// constraint pk_contactValueID
//    pk(id)
type contactValue struct {
	id        int    // integer
	contactID int    // integer
	value     string // varchar(256)
	active    bool
}

func newContactValue() *contactValue {
	o := contactValue{
		id:     0, // seqContactValue.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: person__contact
//
// create table personContactValue
// constraint pk_personContactValue
//    pk(personID)
//    pk(contactValueID)
type personContactValue struct {
	personID       int // integer
	contactValueID int // integer
}

func newPersonContactValue() *personContactValue {
	o := personContactValue{}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  This is the sql that will create the container elements
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the story element table
// create sequence seqStoryElement
// start 1024
type seqStoryElement struct {
	Val int
}

func newSeqStoryElement() *seqStoryElement {
	return &seqStoryElement{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryElement) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryElement) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the media element table
// create sequence seqMediaElement
// start 1024
type seqMediaElement struct {
	Val int
}

func newSeqMediaElement() *seqMediaElement {
	return &seqMediaElement{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaElement) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaElement) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table storyElement
//
//  Description: Holds the properties of a container element. Note that
//               elements can hold either other container elements or field
//               elements, not both.
//
// create table storyElement
// constraint pk_storyElement__id
//    pk(id)
type storyElement struct {
	id               int // integer
	elementTypeID    int // integer
	objectInstanceID int // integer
	parentID         int // integer
	place            int // integer
	objectOrder      int // integer
	displayed        bool
	relatedStoryID   int // integer
	relatedMediaID   int // integer
	active           bool
}

func newStoryElement() *storyElement {
	o := storyElement{
		id:        0, // seqStoryElement.NextVal(),
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
//
// create table mediaElement
// constraint pk_mediaElement__id
//    pk(id)
type mediaElement struct {
	id               int // integer
	elementTypeID    int // integer
	objectInstanceID int // integer
	parentID         int // integer
	place            int // integer
	objectOrder      int // integer
	displayed        bool
	relatedStoryID   int // integer
	relatedMediaID   int // integer
	active           bool
}

func newMediaElement() *mediaElement {
	o := mediaElement{
		id:        0, // seqMediaElement.NextVal(),
		displayed: false,
		active:    true,
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the story field elements.
// create sequence seqStoryField
// start 1024
type seqStoryField struct {
	Val int
}

func newSeqStoryField() *seqStoryField {
	return &seqStoryField{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryField) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryField) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the media field elements.
// create sequence seqMediaField
// start 1024
type seqMediaField struct {
	Val int
}

func newSeqMediaField() *seqMediaField {
	return &seqMediaField{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaField) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaField) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
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
//
// create table storyField
// constraint pk_storyField__id
//    pk(id)
type storyField struct {
	id               int // integer
	fieldTypeID      int // integer
	objectInstanceID int // integer
	parentID         int // integer
	holdVal          bool
	place            int       // integer
	objectOrder      int       // integer
	dateVal          time.Time // timestamp
	shortVal         []byte    // text, like a blob
	blobVal          []byte    // text, like a blob
	active           bool
}

func newStoryField() *storyField {
	o := storyField{
		id:      0, // seqStoryField.NextVal(),
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
//
// create table mediaField
// constraint pk_mediaField__id
//    pk(id)
type mediaField struct {
	id               int // integer
	fieldTypeID      int // integer
	objectInstanceID int // integer
	parentID         int // integer
	place            int // integer
	holdVal          bool
	objectOrder      int       // integer
	dateVal          time.Time // timestamp
	shortVal         string    // varchar(1024)
	blobVal          []byte    // text, like a blob
	active           bool
}

func newMediaField() *mediaField {
	o := mediaField{
		id:      0, // seqMediaField.NextVal(),
		holdVal: false,
		active:  true,
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//
//  This is the SQL that will create the elementType table.
//  It is related to the Bric::ElementType class.
//  Related tables are element and field
//
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the element table
// create sequence seqElementType
// start 1024
type seqElementType struct {
	Val int
}

func newSeqElementType() *seqElementType {
	return &seqElementType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqElementType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqElementType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the subelement table
// create sequence seqSubelementType
// start 1024
type seqSubelementType struct {
	Val int
}

func newSeqSubelementType() *seqSubelementType {
	return &seqSubelementType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqSubelementType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqSubelementType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for element__outputChannel
// create sequence seqElementTypeOutputChannel
// start 1024
type seqElementTypeOutputChannel struct {
	Val int
}

func newSeqElementTypeOutputChannel() *seqElementTypeOutputChannel {
	return &seqElementTypeOutputChannel{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqElementTypeOutputChannel) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqElementTypeOutputChannel) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for element__language
// create sequence seqElementLanguage
// start 1024
type seqElementLanguage struct {
	Val int
}

func newSeqElementLanguage() *seqElementLanguage {
	return &seqElementLanguage{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqElementLanguage) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqElementLanguage) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for elementTypeMember
// create sequence seqElementTypeMember
// start 1024
type seqElementTypeMember struct {
	Val int
}

func newSeqElementTypeMember() *seqElementTypeMember {
	return &seqElementTypeMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqElementTypeMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqElementTypeMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the attr_element table
// create sequence seqAttrElementType
// start 1024
type seqAttrElementType struct {
	Val int
}

func newSeqAttrElementType() *seqAttrElementType {
	return &seqAttrElementType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrElementType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrElementType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for each attr_element_*_val table
// create sequence seqAttrElementTypeVal
// start 1024
type seqAttrElementTypeVal struct {
	Val int
}

func newSeqAttrElementTypeVal() *seqAttrElementTypeVal {
	return &seqAttrElementTypeVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrElementTypeVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrElementTypeVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the element_meta table
// create sequence seqAttrElementTypeMeta
// start 1024
type seqAttrElementTypeMeta struct {
	Val int
}

func newSeqAttrElementTypeMeta() *seqAttrElementTypeMeta {
	return &seqAttrElementTypeMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrElementTypeMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrElementTypeMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the element__site table.
// create sequence seqElementTypeSite
// start 1024
type seqElementTypeSite struct {
	Val int
}

func newSeqElementTypeSite() *seqElementTypeSite {
	return &seqElementTypeSite{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqElementTypeSite) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqElementTypeSite) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: elementType
//
//  Description: The table that holds the information for a given asset type.
//               Holds name and description information and is references by
//               element_contaner and fieldType rows.
//
// create table elementType
// constraint pk_elementTypeID
//    pk(id)
type elementType struct {
	id           int    // integer
	name         string // varchar(64)
	keyName      string // varchar(64)
	description  string // varchar(256)
	topLevel     bool
	paginated    bool
	fixedURI     bool
	relatedStory bool
	relatedMedia bool
	displayed    bool
	media        bool
	bizClassID   int // integer
	typeID       int // integer
	active       bool
}

func newElementType() *elementType {
	o := elementType{
		id:           0, // seqElementType.NextVal(),
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
// create table subelementType
// constraint pk_subelementTypeID
//    pk(id)
type subelementType struct {
	id            int // integer
	parentID      int // integer
	childID       int // integer
	place         int // integer
	minOccurrence int // integer
	maxOccurrence int // integer
}

func newSubelementType() *subelementType {
	o := subelementType{
		id:            0, // seqSubelementType.NextVal(),
		place:         1,
		minOccurrence: 0,
		maxOccurrence: 0,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: element__site
//
//  Description: A table that maps
// create table elementTypeSite
// constraint pk_elementTypeSite__id
//    pk(id)
type elementTypeSite struct {
	id            int // integer
	elementTypeID int // integer
	siteID        int // integer
	active        bool
	primaryOcID   int // integer
}

func newElementTypeSite() *elementTypeSite {
	o := elementTypeSite{
		id:     0, // seqElementTypeSite.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: element__outputChannel
//
//  Description: Holds a reference to the asset type table, the output channel
//               table and an active flag
//
// create table elementTypeOutputChannel
// constraint pk_elementTypeOutputChannel__id
//    pk(id)
type elementTypeOutputChannel struct {
	id              int // integer
	elementTypeID   int // integer
	outputChannelID int // integer
	enabled         bool
	active          bool
}

func newElementTypeOutputChannel() *elementTypeOutputChannel {
	o := elementTypeOutputChannel{
		id:      0, // seqElementTypeOutputChannel.NextVal(),
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
// create table elementTypeMember
// constraint pk_elementTypeMember__id
//    pk(id)
type elementTypeMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newElementTypeMember() *elementTypeMember {
	o := elementTypeMember{
		id: 0, // seqElementTypeMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its element ID and an attribute name.
// create table attrElementType
// constraint pk_attrElementType__id
//    pk(id)
type attrElementType struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrElementType() *attrElementType {
	o := attrElementType{
		id:     0, // seqAttrElementType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element_val
//
//  Description: A table to hold attribute values.
// create table attrElementTypeVal
// constraint pk_attrElementTypeVal__id
//    pk(id)
type attrElementTypeVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrElementTypeVal() *attrElementTypeVal {
	o := attrElementTypeVal{
		id:     0, // seqAttrElementTypeVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attr_element_meta
//
//  Description: A table to represent metadata on types of attributes.
// create table attrElementTypeMeta
// constraint pk_attrElementTypeMeta__id
//    pk(id)
type attrElementTypeMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  string // varchar(2048)
	active bool
}

func newAttrElementTypeMeta() *attrElementTypeMeta {
	o := attrElementTypeMeta{
		id:     0, // seqAttrElementTypeMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  FK index on element__site.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  The sql to create the fieldType table.
//  This maps to the Bric::ElementType::Parts::FieldType class.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the fieldType table
// create sequence seqFieldType
// start 1024
type seqFieldType struct {
	Val int
}

func newSeqFieldType() *seqFieldType {
	return &seqFieldType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqFieldType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqFieldType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the attrFieldType table
// create sequence seqAttrFieldType
// start 1024
type seqAttrFieldType struct {
	Val int
}

func newSeqAttrFieldType() *seqAttrFieldType {
	return &seqAttrFieldType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrFieldType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrFieldType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for each attrFieldTypeVal table
// create sequence seqAttrFieldTypeVal
// start 1024
type seqAttrFieldTypeVal struct {
	Val int
}

func newSeqAttrFieldTypeVal() *seqAttrFieldTypeVal {
	return &seqAttrFieldTypeVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrFieldTypeVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrFieldTypeVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the fieldType_meta table
// create sequence seqAttrFieldTypeMeta
// start 1024
type seqAttrFieldTypeMeta struct {
	Val int
}

func newSeqAttrFieldTypeMeta() *seqAttrFieldTypeMeta {
	return &seqAttrFieldTypeMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrFieldTypeMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrFieldTypeMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: fieldType
//
//  Description: This is the table that contains the name and rules for fields
//          types. It contains references to the elementType table. The place
//          column represents the order that this is to be represented with in
//          its container.
//
// create table fieldType
// constraint pk_fieldTypeID
//    pk(id)
type fieldType struct {
	id            int    // integer
	elementTypeID int    // integer
	name          []byte // text, like a blob
	keyName       []byte // text, like a blob
	description   []byte // text, like a blob
	place         int    // integer
	minOccurrence int    // integer
	maxOccurrence int    // integer
	autopopulated bool
	maxLength     int    // integer
	sqlType       string // varchar(30)
	widgetType    string // varchar(30)
	precision     int    // integer
	cols          int    // integer
	rows          int    // integer
	length        int    // integer
	vals          []byte // text, like a blob
	multiple      bool
	defaultVal    []byte // text, like a blob
	active        bool
}

func newFieldType() *fieldType {
	o := fieldType{
		id:            0, // seqFieldType.NextVal(),
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
// create table attrFieldType
// constraint pk_attrFieldType__id
//    pk(id)
type attrFieldType struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrFieldType() *attrFieldType {
	o := attrFieldType{
		id:     0, // seqAttrFieldType.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrFieldTypeVal
//
//  Description: A table to hold attribute values.
// create table attrFieldTypeVal
// constraint pk_attrFieldTypeVal__id
//    pk(id)
type attrFieldTypeVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrFieldTypeVal() *attrFieldTypeVal {
	o := attrFieldTypeVal{
		id:     0, // seqAttrFieldTypeVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: attrFieldTypeMeta
//
//  Description: A table to represent metadata on types of attributes.
// create table attrFieldTypeMeta
// constraint pk_attrFieldTypeMeta__id
//    pk(id)
type attrFieldTypeMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  []byte // text, like a blob
	active bool
}

func newAttrFieldTypeMeta() *attrFieldTypeMeta {
	o := attrFieldTypeMeta{
		id:     0, // seqAttrFieldTypeMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  This SQL creates the tables necessary for the keyword object.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the main keyword table
// create sequence seqKeyword
// start 1024
type seqKeyword struct {
	Val int
}

func newSeqKeyword() *seqKeyword {
	return &seqKeyword{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqKeyword) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqKeyword) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqKeywordMember
// start 1024
type seqKeywordMember struct {
	Val int
}

func newSeqKeywordMember() *seqKeywordMember {
	return &seqKeywordMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqKeywordMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqKeywordMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: KEYWORD
//
//  Description: The main keyword table.
// create table keyword
// constraint pk_keyword__id
//    pk(id)
type keyword struct {
	id         int    // integer
	name       string // varchar(256)
	screenName string // varchar(256)
	sortName   string // varchar(256)
	active     bool
}

func newKeyword() *keyword {
	o := keyword{
		id:     0, // seqKeyword.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: storyKeyword
//
//  Description: The link between stories and keywords
//
// create table storyKeyword
// constraint pk_storyKeyword
//    pk(storyID)
//    pk(keywordID)
type storyKeyword struct {
	storyID   int // integer
	keywordID int // integer
}

func newStoryKeyword() *storyKeyword {
	o := storyKeyword{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: mediaKeyword
//
//  Description: The link between media and keywords
//
// create table mediaKeyword
// constraint pk_mediaKeyword
//    pk(mediaID)
//    pk(keywordID)
type mediaKeyword struct {
	mediaID   int // integer
	keywordID int // integer
}

func newMediaKeyword() *mediaKeyword {
	o := mediaKeyword{}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: categoryKeyword
//
//  Description: The link between categories and keywords
//
// create table categoryKeyword
// constraint pk_categoryKeyword
//    pk(categoryID)
//    pk(keywordID)
type categoryKeyword struct {
	categoryID int // integer
	keywordID  int // integer
}

func newCategoryKeyword() *categoryKeyword {
	o := categoryKeyword{}
	return &o
}

//
//  TABLE: keywordMember
//
// create table keywordMember
// constraint pk_keywordMember__id
//    pk(id)
type keywordMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newKeywordMember() *keywordMember {
	o := keywordMember{
		id: 0, // seqKeywordMember.NextVal(),
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the table structure for Bric::Org objects.
//
//  SEQUENCES.
//
// create sequence seqOrg
// start 1024
type seqOrg struct {
	Val int
}

func newSeqOrg() *seqOrg {
	return &seqOrg{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqOrg) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqOrg) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: org
//
// create table org
// constraint pk_orgID
//    pk(id)
type org struct {
	id       int    // integer
	name     string // varchar(64)
	longName string // varchar(128)
	personal bool
	active   bool
}

func newOrg() *org {
	o := org{
		id:       0, // seqOrg.NextVal(),
		personal: false,
		active:   true,
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the table structure for Bric::Org::Parts::Address objects.
//
//  SEQUENCES.
//
// create sequence seqAddr
// start 1024
type seqAddr struct {
	Val int
}

func newSeqAddr() *seqAddr {
	return &seqAddr{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAddr) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAddr) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAddrPart
// start 1024
type seqAddrPart struct {
	Val int
}

func newSeqAddrPart() *seqAddrPart {
	return &seqAddrPart{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAddrPart) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAddrPart) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAddrPartType
// start 1024
type seqAddrPartType struct {
	Val int
}

func newSeqAddrPartType() *seqAddrPartType {
	return &seqAddrPartType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAddrPartType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAddrPartType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: addr
//
// create table addr
// constraint pk_addrID
//    pk(id)
type addr struct {
	id     int    // integer
	orgID  int    // integer
	ttypee string // varchar(64)
	active bool
}

func newAddr() *addr {
	o := addr{
		id:     0, // seqAddr.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: addrPartType
//
// create table addrPartType
// constraint pk_addrPartTypeID
//    pk(id)
type addrPartType struct {
	id     int    // integer
	name   string // varchar(64)
	active bool
}

func newAddrPartType() *addrPartType {
	o := addrPartType{
		id:     0, // seqAddrPartType.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: addrPart
//
// create table addrPart
// constraint pk_addrPart__id
//    pk(id)
type addrPart struct {
	id             int    // integer
	addrID         int    // integer
	addrPartTypeID int    // integer
	value          string // varchar(256)
}

func newAddrPart() *addrPart {
	o := addrPart{
		id: 0, // seqAddrPart.NextVal(),
	}
	return &o
}

//
//  TABLE: personOrgAddr
//
// create table personOrgAddr
// constraint pk_personOrgAddr__all
//    pk(addrID)
//    pk(personOrgID)
type personOrgAddr struct {
	addrID      int // integer
	personOrgID int // integer
}

func newPersonOrgAddr() *personOrgAddr {
	o := personOrgAddr{}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the table structure for Bric::Org::Person objects.
//
//  SEQUENCES.
//
// create sequence seqPersonOrg
// start 1024
type seqPersonOrg struct {
	Val int
}

func newSeqPersonOrg() *seqPersonOrg {
	return &seqPersonOrg{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPersonOrg) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPersonOrg) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: personOrg
//
// create table personOrg
// constraint pk_personOrgID
//    pk(id)
type personOrg struct {
	id         int    // integer
	personID   int    // integer
	orgID      int    // integer
	role       string // varchar(64)
	department string // varchar(64)
	title      string // varchar(64)
	active     bool
}

func newPersonOrg() *personOrg {
	o := personOrg{
		id:     0, // seqPersonOrg.NextVal(),
		active: true,
	}
	return &o
}

//
//  INDEXES.
//
//
//  Project: Bricolage Business API
//
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the table structure for Bric::BC::Org::Source objects.
//
//  SEQUENCES.
//
// create sequence seqSource
// start 1024
type seqSource struct {
	Val int
}

func newSeqSource() *seqSource {
	return &seqSource{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqSource) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqSource) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: source
//
// create table source
// constraint pk_sourceID
//    pk(id)
type source struct {
	id          int    // integer
	orgID       int    // integer
	name        string // varchar(64)
	description string // varchar(256)
	expire      int    // integer
	active      bool
}

func newSource() *source {
	o := source{
		id:     0, // seqSource.NextVal(),
		expire: 0,
		active: true,
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  Description: The table that holds the registered Output Channels.
//                 This maps to the Bric::OutputChannel Class.
//
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the outputChannel table
// create sequence seqOutputChannel
// start 1024
type seqOutputChannel struct {
	Val int
}

func newSeqOutputChannel() *seqOutputChannel {
	return &seqOutputChannel{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqOutputChannel) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqOutputChannel) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqOutputChannelInclude
// start 1024
type seqOutputChannelInclude struct {
	Val int
}

func newSeqOutputChannelInclude() *seqOutputChannelInclude {
	return &seqOutputChannelInclude{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqOutputChannelInclude) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqOutputChannelInclude) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqOutputChannelMember
// start 1024
type seqOutputChannelMember struct {
	Val int
}

func newSeqOutputChannelMember() *seqOutputChannelMember {
	return &seqOutputChannelMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqOutputChannelMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqOutputChannelMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table outputChannel
//
//  Description: Holds info on the various output channels and is referenced
//                  by templates and elements
//
//
// create table outputChannel
// constraint pk_outputChannelID
//    pk(id)
type outputChannel struct {
	id             int    // integer
	name           string // varchar(64)
	description    string // varchar(256)
	siteID         int    // integer
	protocol       string // varchar(16)
	filename       string // varchar(32)
	fileExt        string // varchar(32)
	primaryCe      bool
	uriFormat      string // varchar(64)
	fixedURIFormat string // varchar(64)
	uriCase        int    // integer
	useSlug        bool
	burner         int // integer
	active         bool
}

func newOutputChannel() *outputChannel {
	o := outputChannel{
		id:      0, // seqOutputChannel.NextVal(),
		uriCase: 1,
		useSlug: false,
		burner:  1,
		active:  true,
	}
	// check uriCase ( uriCase  in  1  ,  2  ,  3 )
	return &o
}

//
//  TABLE: outputChannelInclude
//
// create table outputChannelInclude
// constraint pk_outputChannelInclude__id
//    pk(id)
type outputChannelInclude struct {
	id              int // integer
	outputChannelID int // integer
	includeOcID     int // integer
}

func newOutputChannelInclude() *outputChannelInclude {
	o := outputChannelInclude{
		id: 0, // seqOutputChannelInclude.NextVal(),
	}
	// check includeOcID ( includeOcID  <  >  outputChannelID )
	return &o
}

//
//  TABLE: outputChannelMember
//
// create table outputChannelMember
// constraint pk_outputChannelMember__id
//    pk(id)
type outputChannelMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newOutputChannelMember() *outputChannelMember {
	o := outputChannelMember{
		id: 0, // seqOutputChannelMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the basic table for all Bric::Person objects. The indexes are
//  suggestions.
//
//  SEQUENCES.
//
// create sequence seqPerson
// start 1024
type seqPerson struct {
	Val int
}

func newSeqPerson() *seqPerson {
	return &seqPerson{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPerson) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPerson) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqPersonMember
// start 1024
type seqPersonMember struct {
	Val int
}

func newSeqPersonMember() *seqPersonMember {
	return &seqPersonMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPersonMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPersonMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: person
//
// create table person
// constraint pk_personID
//    pk(id)
type person struct {
	id     int    // integer
	prefix string // varchar(32)
	lname  string // varchar(64)
	fname  string // varchar(64)
	mname  string // varchar(64)
	suffix string // varchar(32)
	active bool
}

func newPerson() *person {
	o := person{
		id:     0, // seqPerson.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: personMember
//
// create table personMember
// constraint pk_personMember__id
//    pk(id)
type personMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newPersonMember() *personMember {
	o := personMember{
		id: 0, // seqPersonMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the basic table for Bric::Person::Usr objects, and
//  establishes its relationship with Bric::Person. The login field must be unique,
//  hence the udx_usr__login index.
//
//  TABLE: usr
//
// create table usr
// constraint pk_usrID
//    pk(id)
type usr struct {
	id       int    // integer
	login    string // varchar(128)
	password string // char(32)
	active   bool
}

func newUsr() *usr {
	o := usr{
		active: true,
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the basic table for all Bric::Site objects. The indexes are
//  suggestions.
//
//  SEQUENCES.
//
// create sequence seqSiteMember
// start 1024
type seqSiteMember struct {
	Val int
}

func newSeqSiteMember() *seqSiteMember {
	return &seqSiteMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqSiteMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqSiteMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: site
//
// create table site
// constraint pk_siteID
//    pk(id)
type site struct {
	id          int    // integer
	name        []byte // text, like a blob
	description []byte // text, like a blob
	domainName  []byte // text, like a blob
	active      bool
}

func newSite() *site {
	o := site{
		active: true,
	}
	return &o
}

//
//  TABLE: siteMember
//
// create table siteMember
// constraint pk_siteMember__id
//    pk(id)
type siteMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newSiteMember() *siteMember {
	o := siteMember{
		id: 0, // seqSiteMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  The database tables for the Workflow class.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  A sequence of unique IDs for the workflow table.
// create sequence seqWorkflow
// start 1024
type seqWorkflow struct {
	Val int
}

func newSeqWorkflow() *seqWorkflow {
	return &seqWorkflow{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqWorkflow) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqWorkflow) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqWorkflowMember
// start 1024
type seqWorkflowMember struct {
	Val int
}

func newSeqWorkflowMember() *seqWorkflowMember {
	return &seqWorkflowMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqWorkflowMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqWorkflowMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: workflow
//
//  Description: The main workflow table.
// create table workflow
// constraint pk_workflowID
//    pk(id)
type workflow struct {
	id           int    // integer
	name         string // varchar(64)
	description  string // varchar(256)
	allDeskGrpID int    // integer
	reqDeskGrpID int    // integer
	assetGrpID   int    // integer
	headDeskID   int    // integer
	ttypee       int    // integer
	active       bool
	siteID       int // integer
}

func newWorkflow() *workflow {
	o := workflow{
		id:     0, // seqWorkflow.NextVal(),
		ttypee: 1,
		active: true,
	}
	// check type ( type  in  1  ,  2  ,  3 )
	return &o
}

//
//  TABLE: workflowMember
//
// create table workflowMember
// constraint pk_workflowMember__id
//    pk(id)
type workflowMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newWorkflowMember() *workflowMember {
	o := workflowMember{
		id: 0, // seqWorkflowMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  The database tables for the Bric::Workflow::Parts::Desk class.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  A sequence of unique IDs for the desk table.
// create sequence seqDesk
// start 1024
type seqDesk struct {
	Val int
}

func newSeqDesk() *seqDesk {
	return &seqDesk{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqDesk) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqDesk) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the desk member ordering.
// create sequence seqDeskMember
// start 1024
type seqDeskMember struct {
	Val int
}

func newSeqDeskMember() *seqDeskMember {
	return &seqDeskMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqDeskMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqDeskMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: desk
//
//  Description: Represents a desk in the workflow
// create table desk
// constraint pk_deskID
//    pk(id)
type desk struct {
	id           int    // integer
	name         string // varchar(64)
	description  string // varchar(256)
	preChkRules  int    // integer
	postChkRules int    // integer
	assetGrp     int    // integer
	publish      bool
	active       bool
}

func newDesk() *desk {
	o := desk{
		id:      0, // seqDesk.NextVal(),
		publish: false,
		active:  true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Table: deskMember
//
//  Description: The link between desk objects and member objects
//
// create table deskMember
// constraint pk_deskMember__id
//    pk(id)
type deskMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newDeskMember() *deskMember {
	o := deskMember{
		id: 0, // seqDeskMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAction
// start 1024
type seqAction struct {
	Val int
}

func newSeqAction() *seqAction {
	return &seqAction{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAction) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAction) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAttrAction
// start 1024
type seqAttrAction struct {
	Val int
}

func newSeqAttrAction() *seqAttrAction {
	return &seqAttrAction{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrAction) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrAction) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAttrActionVal
// start 1024
type seqAttrActionVal struct {
	Val int
}

func newSeqAttrActionVal() *seqAttrActionVal {
	return &seqAttrActionVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrActionVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrActionVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqAttrActionMeta
// start 1024
type seqAttrActionMeta struct {
	Val int
}

func newSeqAttrActionMeta() *seqAttrActionMeta {
	return &seqAttrActionMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrActionMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrActionMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: action
//
// create table action
// constraint pk_action__id
//    pk(id)
type action struct {
	id           int // integer
	ord          int // integer
	serverTypeID int // integer
	actionTypeID int // integer
	active       bool
}

func newAction() *action {
	o := action{
		id:     0, // seqAction.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: attrAction
//
// create table attrAction
// constraint pk_attrAction__id
//    pk(id)
type attrAction struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrAction() *attrAction {
	o := attrAction{
		id:     0, // seqAttrAction.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: attrActionVal
//
// create table attrActionVal
// constraint pk_attrActionVal__id
//    pk(id)
type attrActionVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrActionVal() *attrActionVal {
	o := attrActionVal{
		id:     0, // seqAttrActionVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//
//  TABLE: attrActionMeta
//
// create table attrActionMeta
// constraint pk_attrActionMeta__id
//    pk(id)
type attrActionMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  string // varchar(2048)
	active bool
}

func newAttrActionMeta() *attrActionMeta {
	o := attrActionMeta{
		id:     0, // seqAttrActionMeta.NextVal(),
		active: true,
	}
	return &o
}

//
//  INDEXES.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqActionType
// start 1024
type seqActionType struct {
	Val int
}

func newSeqActionType() *seqActionType {
	return &seqActionType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqActionType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqActionType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: actionType
//
// create table actionType
// constraint pk_actionTypeID
//    pk(id)
type actionType struct {
	id          int    // integer
	name        string // varchar(64)
	description string // varchar(256)
	active      bool
}

func newActionType() *actionType {
	o := actionType{
		id:     0, // seqActionType.NextVal(),
		active: false,
	}
	return &o
}

//
//  TABLE: actionTypeMediaType
//
// create table actionTypeMediaType
// constraint pk_action__mediaType
//    pk(actionTypeID)
//    pk(mediaTypeID)
type actionTypeMediaType struct {
	actionTypeID int // integer
	mediaTypeID  int // integer
}

func newActionTypeMediaType() *actionTypeMediaType {
	o := actionTypeMediaType{}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage Business API
//  File:    Resource.sql
//
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqResource
// start 1024
type seqResource struct {
	Val int
}

func newSeqResource() *seqResource {
	return &seqResource{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqResource) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqResource) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: mediaResource
//
// create table mediaResource
// constraint pk_mediaResource
//    pk(mediaID)
//    pk(resourceID)
type mediaResource struct {
	resourceID int // integer
	mediaID    int // integer
}

func newMediaResource() *mediaResource {
	o := mediaResource{}
	return &o
}

//
//  TABLE: resource
//
// create table resource
// constraint pk_resourceID
//    pk(id)
type resource struct {
	id          int       // integer
	parentID    int       // integer
	mediaTypeID int       // integer
	path        string    // varchar(256)
	uri         string    // varchar(256)
	size        int       // integer
	modTime     time.Time // timestamp
	isDir       bool
}

func newResource() *resource {
	o := resource{
		id: 0, // seqResource.NextVal(),
	}
	return &o
}

//
//  TABLE: storyResource
//
// create table storyResource
// constraint pk_storyResource
//    pk(storyID)
//    pk(resourceID)
type storyResource struct {
	storyID    int // integer
	resourceID int // integer
}

func newStoryResource() *storyResource {
	o := storyResource{}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  Sequences.
//
// create sequence seqServer
// start 1024
type seqServer struct {
	Val int
}

func newSeqServer() *seqServer {
	return &seqServer{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqServer) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqServer) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: server
//
// create table server
// constraint pk_server__id
//    pk(id)
type server struct {
	id           int    // integer
	serverTypeID int    // integer
	hostName     string // varchar(128)
	os           string // char(5)
	docRoot      string // varchar(128)
	login        string // varchar(64)
	password     string // varchar(64)
	cookie       string // varchar(512)
	active       bool
}

func newServer() *server {
	o := server{
		id:     0, // seqServer.NextVal(),
		active: true,
	}
	return &o
}

//
//  Indexes.
//
//  Project: Bricolage Business API
//  File:    ServerType.sql
//
//  Author: David Wheeler <david@justatheory.com>
//
//
//  Sequences.
//
// create sequence seqServerType
// start 1024
type seqServerType struct {
	Val int
}

func newSeqServerType() *seqServerType {
	return &seqServerType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqServerType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqServerType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqDestMember
// start 1024
type seqDestMember struct {
	Val int
}

func newSeqDestMember() *seqDestMember {
	return &seqDestMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqDestMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqDestMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: serverType
//
// create table serverType
// constraint pk_serverTypeID
//    pk(id)
type serverType struct {
	id          int    // integer
	classID     int    // integer
	name        string // varchar(64)
	description string // varchar(256)
	siteID      int    // integer
	copyable    bool
	publish     bool
	preview     bool
	active      bool
}

func newServerType() *serverType {
	o := serverType{
		id:       0, // seqServerType.NextVal(),
		copyable: false,
		publish:  true,
		preview:  false,
		active:   true,
	}
	return &o
}

//
//  TABLE: serverType__output__channel
//
// create table serverTypeOutputChannel
// constraint pk_serverTypeOutputChannel
//    pk(serverTypeID)
//    pk(outputChannelID)
type serverTypeOutputChannel struct {
	serverTypeID    int // integer
	outputChannelID int // integer
}

func newServerTypeOutputChannel() *serverTypeOutputChannel {
	o := serverTypeOutputChannel{}
	return &o
}

//
//  TABLE: destMember
//
// create table destMember
// constraint pk_destMember__id
//    pk(id)
type destMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newDestMember() *destMember {
	o := destMember{
		id: 0, // seqDestMember.NextVal(),
	}
	return &o
}

//
//  Indexes.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAlert
// start 1024
type seqAlert struct {
	Val int
}

func newSeqAlert() *seqAlert {
	return &seqAlert{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAlert) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAlert) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: alert
//
// create table alert
// constraint pk_alertID
//    pk(id)
type alert struct {
	id          int       // integer
	alertTypeID int       // integer
	eventID     int       // integer
	subject     string    // varchar(128)
	message     string    // varchar(512)
	timestamp   time.Time // timestamp
}

func newAlert() *alert {
	o := alert{
		id:        0, // seqAlert.NextVal(),
		timestamp: time.Now(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAlertType
// start 1024
type seqAlertType struct {
	Val int
}

func newSeqAlertType() *seqAlertType {
	return &seqAlertType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAlertType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAlertType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: alertType
//
// create table alertType
// constraint pk_alertTypeID
//    pk(id)
type alertType struct {
	id          int    // integer
	eventTypeID int    // integer
	usrID       int    // integer
	name        string // varchar(64)
	subject     string // varchar(128)
	message     string // varchar(512)
	active      bool
	del         bool
}

func newAlertType() *alertType {
	o := alertType{
		id:     0, // seqAlertType.NextVal(),
		active: true,
		del:    false,
	}
	return &o
}

//
//  TABLE: alertTypeGrpContact
//
// create table alertTypeGrpContact
// constraint pk_alertTypeGrpContact
//    pk(alertTypeID)
//    pk(contactID)
//    pk(grpID)
type alertTypeGrpContact struct {
	alertTypeID int // integer
	contactID   int // integer
	grpID       int // integer
}

func newAlertTypeGrpContact() *alertTypeGrpContact {
	o := alertTypeGrpContact{}
	return &o
}

//
//  TABLE: alertTypeUsrContact
//
// create table alertTypeUsrContact
// constraint pk_alertTypeUsrContact
//    pk(alertTypeID)
//    pk(usrID)
//    pk(contactID)
type alertTypeUsrContact struct {
	alertTypeID int // integer
	contactID   int // integer
	usrID       int // integer
}

func newAlertTypeUsrContact() *alertTypeUsrContact {
	o := alertTypeUsrContact{}
	return &o
}

//
//  INDEXES.
//
//  alertType
//  alertTypeGrpContact
//  alertTypeUsrContact
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAlertTypeRule
// start 1024
type seqAlertTypeRule struct {
	Val int
}

func newSeqAlertTypeRule() *seqAlertTypeRule {
	return &seqAlertTypeRule{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAlertTypeRule) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAlertTypeRule) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: alertTypeRule
//
// create table alertTypeRule
// constraint pk_alertTypeRule__id
//    pk(id)
type alertTypeRule struct {
	id          int    // integer
	alertTypeID int    // integer
	attr        string // varchar(64)
	operator    string // char(2)
	value       string // varchar(256)
}

func newAlertTypeRule() *alertTypeRule {
	o := alertTypeRule{
		id: 0, // seqAlertTypeRule.NextVal(),
	}
	return &o
}

//
//  INDEXS.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAlerted
// start 1024
type seqAlerted struct {
	Val int
}

func newSeqAlerted() *seqAlerted {
	return &seqAlerted{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAlerted) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAlerted) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: alerted
//
// create table alerted
// constraint pk_alertedID
//    pk(id)
type alerted struct {
	id      int       // integer
	usrID   int       // integer
	alertID int       // integer
	ackTime time.Time // timestamp
}

func newAlerted() *alerted {
	o := alerted{
		id: 0, // seqAlerted.NextVal(),
	}
	return &o
}

//
//  TABLE: alertedContactValue
//
// create table alertedContactValue
// constraint pk_alertedContactValue
//    pk(alertedID)
//    pk(contactID)
//    pk(contactValueValue)
type alertedContactValue struct {
	alertedID         int       // integer
	contactID         int       // integer
	contactValueValue string    // varchar(256)
	sentTime          time.Time // timestamp
}

func newAlertedContactValue() *alertedContactValue {
	o := alertedContactValue{}
	return &o
}

//
//  INDEXES.
//
//  alerted
//  alertedContactValue
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  -----------------------------------------------------------------------------
//  Attribute.sql
//
//
//  This SQL creates the tables necessary for the attribute object.  This file
//  applies to attributes on the Bric::Person class.  Any other classes that
//  require attributes need only duplicate these tables, changing 'person' to
//  the correct class name.  Class names may be shortened to ensure that the
//  resulting table names are under the oracle 30 character name limit as long
//  as the resulting shortened class name is unique.
//
// 2594: /* Commented out because attr_person won't be used in production (in this version).
// 2594:    However, the examples still apply. --David, 23 Feb 2001
// 2594:
// 2594: -- -----------------------------------------------------------------------------
// 2594: -- Sequences
// 2594:
// 2594: -- Unique IDs for the attr_person table
// 2594: CREATE SEQUENCE seq_attr_person START 1024;
// 2594:
// 2594: -- Unique IDs for each attr_person_*_val table
// 2594: CREATE SEQUENCE seq_attr_person_val START 1024;
// 2594:
// 2594: -- Unique IDs for the person_meta table
// 2594: CREATE SEQUENCE seq_attr_person_meta START 1024;
// 2594:
// 2594: -- -----------------------------------------------------------------------------
// 2594: -- Table: attr_person
// 2594: --
// 2594: -- Description: A table to represent types of attributes.  A type is defined by
// 2594: --              its subsystem, its person ID and an attribute name.
// 2594:
// 2594: CREATE TABLE attr_person (
// 2594:     id         INTEGER       NOT NULL
// 2594:                              DEFAULT NEXTVAL('seq_attr_person'),
// 2594:     subsys     VARCHAR(256)  NOT NULL,
// 2594:     name       VARCHAR(256)  NOT NULL,
// 2594:     sqlType   VARCHAR(30)   NOT NULL,
// 2594:     active     BOOLEAN       NOT NULL DEFAULT TRUE,
// 2594:    CONSTRAINT pk_attr_personID PRIMARY KEY (id)
// 2594: );
// 2594:
// 2594:
// 2594: -- -----------------------------------------------------------------------------
// 2594: -- Table: attr_person_val
// 2594: --
// 2594: -- Description: A table to hold attribute values.
// 2594:
// 2594: CREATE TABLE attr_person_val (
// 2594:     id           INTEGER         NOT NULL
// 2594:                                  DEFAULT NEXTVAL('seq_attr_person_val'),
// 2594:     objectID   INTEGER         NOT NULL,
// 2594:     attrID     INTEGER         NOT NULL,
// 2594:     dateVal     TIMESTAMP,
// 2594:     shortVal    VARCHAR(1024),
// 2594:     blobVal     TEXT,
// 2594:     serial       BOOLEAN         DEFAULT FALSE,
// 2594:     active       BOOLEAN         NOT NULL DEFAULT TRUE,
// 2594:     CONSTRAINT pk_attr_person_val__id PRIMARY KEY (id)
// 2594: );
// 2594:
// 2594: -- -----------------------------------------------------------------------------
// 2594: -- Table: attr_person_meta
// 2594: --
// 2594: -- Description: A table to represent metadata on types of attributes.
// 2594:
// 2594: CREATE TABLE attr_person_meta (
// 2594:     id        INTEGER         NOT NULL
// 2594:                               DEFAULT NEXTVAL('seq_attr_person_meta'),
// 2594:     attrID  INTEGER         NOT NULL,
// 2594:     name      VARCHAR(256)    NOT NULL,
// 2594:     value     VARCHAR(2048),
// 2594:     active    BOOLEAN         NOT NULL DEFAULT TRUE,
// 2594:    CONSTRAINT pk_attr_person_meta__id PRIMARY KEY (id)
// 2594: );
// 2594:
// 2594:
// 2594: -- -----------------------------------------------------------------------------
// 2594: -- Indexes.
// 2594: --
// 2594:
// 2594: -- Unique index on subsystem/name pair
// 2594: CREATE UNIQUE INDEX udx_attr_person__subsys__name ON attr_person(subsys, name);
// 2594:
// 2594: -- Indexes on name and subsys.
// 2594: CREATE INDEX idx_attr_person__name ON attr_person(LOWER(name));
// 2594: CREATE INDEX idx_attr_person__subsys ON attr_person(LOWER(subsys));
// 2594:
// 2594: -- Unique index on objectID/attrID pair
// 2594: CREATE UNIQUE INDEX udx_attr_person_val__obj_attr ON attr_person_val (objectID,attrID);
// 2594:
// 2594: -- FK indexes on objectID and attrID.
// 2594: CREATE INDEX fkx_person__attr_person_val ON attr_person_val(objectID);
// 2594: CREATE INDEX fkx_attr_person__attr_person_val ON attr_person_val(attrID);
// 2594:
// 2594: -- Unique index on attrID/name pair
// 2594: CREATE UNIQUE INDEX udx_attr_person_meta__attr_name ON attr_person_meta (attrID, name);
// 2594:
// 2594: -- Index on meta name.
// 2594: CREATE INDEX idx_attr_person_meta__name ON attr_person_meta(LOWER(name));
// 2594:
// 2594: -- FK index on attrID.
// 2594: CREATE INDEX fkx_attr_person__attr_person_meta ON attr_person_meta(attrID);
// 2594:
// 2594: */
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Garth Webb <garth@perijove.com>
//
//  This is the SQL that will create the class table.
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the class table
// create sequence seqClass
// start 1024
type seqClass struct {
	Val int
}

func newSeqClass() *seqClass {
	return &seqClass{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqClass) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqClass) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  TABLE: class
//         For keeping track of Perl classes.
// create table class
// constraint pk_classID
//    pk(id)
type class struct {
	id          int    // integer
	keyName     string // varchar(32)
	pkgName     string // varchar(128)
	dispName    string // varchar(128)
	pluralName  string // varchar(128)
	description string // varchar(256)
	distributor bool
}

func newClass() *class {
	o := class{
		id:          0, // seqClass.NextVal(),
		distributor: false,
	}
	// check keyName ( lower  keyName  =  keyName )
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the basic tables for Bric::Util::Event objects.
//
//  SEQUENCES.
//
// create sequence seqEvent
// start 1024
type seqEvent struct {
	Val int
}

func newSeqEvent() *seqEvent {
	return &seqEvent{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqEvent) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqEvent) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqEventAttr
// start 1024
type seqEventAttr struct {
	Val int
}

func newSeqEventAttr() *seqEventAttr {
	return &seqEventAttr{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqEventAttr) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqEventAttr) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: event
//
// create table event
// constraint pk_eventID
//    pk(id)
type event struct {
	id          int       // integer
	eventTypeID int       // integer
	usrID       int       // integer
	objID       int       // integer
	timestamp   time.Time // timestamp
}

func newEvent() *event {
	o := event{
		id:        0, // seqEvent.NextVal(),
		timestamp: time.Now(),
	}
	return &o
}

//
//  TABLE: eventAttr
//
// create table eventAttr
// constraint pk_eventAttr__id
//    pk(id)
type eventAttr struct {
	id              int    // integer
	eventID         int    // integer
	eventTypeAttrID int    // integer
	value           string // varchar(128)
}

func newEventAttr() *eventAttr {
	o := eventAttr{
		id: 0, // seqEventAttr.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//
//  ER/Studio 4.0 SQL Code Generation
//  Project:      Bricolage Business API
//
//  Target DBMS : Oracle 8
//  Author: David Wheeler <david@justatheory.com>
//  This DDL creates the basic table for all Bric::Util::EventType objects. It's
//  pretty easy - they're really just all groups.
//
//  SEQUENCES.
//
// create sequence seqEventType
// start 1024
type seqEventType struct {
	Val int
}

func newSeqEventType() *seqEventType {
	return &seqEventType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqEventType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqEventType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqEventTypeAttr
// start 1024
type seqEventTypeAttr struct {
	Val int
}

func newSeqEventTypeAttr() *seqEventTypeAttr {
	return &seqEventTypeAttr{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqEventTypeAttr) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqEventTypeAttr) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: eventType
//
// create table eventType
// constraint pk_eventTypeID
//    pk(id)
type eventType struct {
	id          int    // integer
	keyName     string // varchar(64)
	name        string // varchar(64)
	description string // varchar(256)
	classID     int    // integer
	active      bool
}

func newEventType() *eventType {
	o := eventType{
		id:     0, // seqEventType.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: eventTypeAttr
//
// create table eventTypeAttr
// constraint pk_eventTypeAttrID
//    pk(id)
type eventTypeAttr struct {
	id          int    // integer
	eventTypeID int    // integer
	name        string // varchar(64)
}

func newEventTypeAttr() *eventTypeAttr {
	o := eventTypeAttr{
		id: 0, // seqEventTypeAttr.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  ----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the grp table
// create sequence seqGrp
// start 1024
type seqGrp struct {
	Val int
}

func newSeqGrp() *seqGrp {
	return &seqGrp{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqGrp) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqGrp) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  ----------------------------------------------------------------------------
//  Table grp
//
//  Description: The grp table   Contains the name and description of the
//                  group and its parent if it has one
//
// create table grp
// constraint pk_grpID
//    pk(id)
type grp struct {
	id          int    // integer
	parentID    int    // integer
	classID     int    // integer
	name        string // varchar(64)
	description string // varchar(256)
	secret      bool
	permanent   bool
	active      bool
}

func newGrp() *grp {
	o := grp{
		id:        0, // seqGrp.NextVal(),
		secret:    true,
		permanent: false,
		active:    true,
	}
	// check parentID ( parentID  <  >  id )
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqAlertTypeMember
// start 1024
type seqAlertTypeMember struct {
	Val int
}

func newSeqAlertTypeMember() *seqAlertTypeMember {
	return &seqAlertTypeMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAlertTypeMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAlertTypeMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: alertTypeMember
//
// create table alertTypeMember
// constraint pk_alertTypeMember__id
//    pk(id)
type alertTypeMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newAlertTypeMember() *alertTypeMember {
	o := alertTypeMember{
		id: 0, // seqAlertTypeMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
// create sequence seqContribTypeMember
// start 1024
type seqContribTypeMember struct {
	Val int
}

func newSeqContribTypeMember() *seqContribTypeMember {
	return &seqContribTypeMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqContribTypeMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqContribTypeMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: contribTypeMember
//
// create table contribTypeMember
// constraint pk_contribTypeMember__id
//    pk(id)
type contribTypeMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newContribTypeMember() *contribTypeMember {
	o := contribTypeMember{
		id: 0, // seqContribTypeMember.NextVal(),
	}
	return &o
}

//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//  TABLE: event_member
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqOrgMember
// start 1024
type seqOrgMember struct {
	Val int
}

func newSeqOrgMember() *seqOrgMember {
	return &seqOrgMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqOrgMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqOrgMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: orgMember
//
// create table orgMember
// constraint pk_orgMember__id
//    pk(id)
type orgMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newOrgMember() *orgMember {
	o := orgMember{
		id: 0, // seqOrgMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  -----------------------------------------------------------------------------
//  Member.sql
//
//
//  The member table and the tables that map member back to their respective
//  objects. The member table contains an id and a group id. The table that
//  maps the object to its member contains an id an object id and a member id
//
//  Thought should be given to:
//          Ensuring that an object is not placed with in the same group twice
//         Making sure that an object that is deactivated from a group that is
//             then put back in again will behave properly
//
//  -----------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the grpMember table
// create sequence seqGrpMember
// start 1024
type seqGrpMember struct {
	Val int
}

func newSeqGrpMember() *seqGrpMember {
	return &seqGrpMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqGrpMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqGrpMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: grpMember
//
//  Description: The link between stroy objects and member objects
//
// create table grpMember
// constraint pk_grpMemberID
//    pk(id)
type grpMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newGrpMember() *grpMember {
	o := grpMember{
		id: 0, // seqGrpMember.NextVal(),
	}
	return &o
}

//
//  INDEXES
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  -----------------------------------------------------------------------------
//  Member.sql
//
//
//  The member table and the tables that map member bac to their respective
//  objects.   The member table contains an id and a group id.   The table that
//  maps the object to its member contains an id an object id and a member id
//
//  Thought should be given to:
//          Ensuring that an object is not placed with in the same group twice
//         Making sure that an object that is deactivated from a group that is
//             then put back in again will behave properly
//
//  -----------------------------------------------------------------------------
//
//  SEQUENCES.
//
// create sequence seqMember
// start 1024
type seqMember struct {
	Val int
}

func newSeqMember() *seqMember {
	return &seqMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqStoryMember
// start 1024
type seqStoryMember struct {
	Val int
}

func newSeqStoryMember() *seqStoryMember {
	return &seqStoryMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqStoryMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqStoryMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  -----------------------------------------------------------------------------
//  Table: member
//
//  Description:    The table that creates a member of a group.   The obj_member
//  table then links the objects to the member table
//
// create table member
// constraint pk_memberID
//    pk(id)
type member struct {
	id      int // integer
	grpID   int // integer
	classID int // integer
	active  bool
}

func newMember() *member {
	o := member{
		id:     0, // seqMember.NextVal(),
		active: true,
	}
	return &o
}

//
//  INDEXES.
//
//  Use the below section as an example to create new member tables for
//  other objects.
//  -----------------------------------------------------------------------------
//  Table: storyMember
//
//  Description: The link between story objects and member objects
//
// create table storyMember
// constraint pk_storyMember__id
//    pk(id)
type storyMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newStoryMember() *storyMember {
	o := storyMember{
		id: 0, // seqStoryMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  -----------------------------------------------------------------------------
//
//  This SQL creates the tables necessary for the attribute object.  This file
//  applies to attributes on the Bric::Util::Grp class.  Any other classes that
//  require attributes need only duplicate these tables, changing 'member' to
//  the correct class name.  Class names may be shortened to ensure that the
//  resulting table names are under the oracle 30 character name limit as long
//  as the resulting shortened class name is unique.
//
//  -------------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the attrMember table
// create sequence seqAttrMember
// start 1024
type seqAttrMember struct {
	Val int
}

func newSeqAttrMember() *seqAttrMember {
	return &seqAttrMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for each attrMember_*_val table
// create sequence seqAttrMemberVal
// start 1024
type seqAttrMemberVal struct {
	Val int
}

func newSeqAttrMemberVal() *seqAttrMemberVal {
	return &seqAttrMemberVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrMemberVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrMemberVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the member_meta table
// create sequence seqAttrMemberMeta
// start 1024
type seqAttrMemberMeta struct {
	Val int
}

func newSeqAttrMemberMeta() *seqAttrMemberMeta {
	return &seqAttrMemberMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrMemberMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrMemberMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Table: attrMember
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its member ID and an attribute name.
// create table attrMember
// constraint pk_attrMember__id
//    pk(id)
type attrMember struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrMember() *attrMember {
	o := attrMember{
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrMemberVal
//  Description: A table to hold attribute values.
// create table attrMemberVal
// constraint pk_attrMemberVal__id
//    pk(id)
type attrMemberVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrMemberVal() *attrMemberVal {
	o := attrMemberVal{
		serial: false,
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrMemberMeta
//  Description: A table to represent metadata on types of attributes.
// create table attrMemberMeta
// constraint pk_attrMemberMeta__id
//    pk(id)
type attrMemberMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  string // varchar(2048)
	active bool
}

func newAttrMemberMeta() *attrMemberMeta {
	o := attrMemberMeta{
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqSourceMember
// start 1024
type seqSourceMember struct {
	Val int
}

func newSeqSourceMember() *seqSourceMember {
	return &seqSourceMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqSourceMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqSourceMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: sourceMember
//
// create table sourceMember
// constraint pk_sourceMember__id
//    pk(id)
type sourceMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newSourceMember() *sourceMember {
	o := sourceMember{
		id: 0, // seqSourceMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqUserMember
// start 1024
type seqUserMember struct {
	Val int
}

func newSeqUserMember() *seqUserMember {
	return &seqUserMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqUserMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqUserMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: userMember
//
// create table userMember
// constraint pk_userMember__id
//    pk(id)
type userMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newUserMember() *userMember {
	o := userMember{
		id: 0, // seqUserMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
//  -----------------------------------------------------------------------------
//
//  This SQL creates the tables necessary for the attribute object.  This file
//  applies to attributes on the Bric::Util::Grp class.  Any other classes that
//  require attributes need only duplicate these tables, changing 'grp' to
//  the correct class name.  Class names may be shortened to ensure that the
//  resulting table names are under the oracle 30 character name limit as long
//  as the resulting shortened class name is unique.
//
//  -------------------------------------------------------------------------------
//  Sequences
//  Unique IDs for the attrGrp table
// create sequence seqAttrGrp
// start 1024
type seqAttrGrp struct {
	Val int
}

func newSeqAttrGrp() *seqAttrGrp {
	return &seqAttrGrp{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrGrp) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrGrp) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for each attrGrp_*_val table
// create sequence seqAttrGrpVal
// start 1024
type seqAttrGrpVal struct {
	Val int
}

func newSeqAttrGrpVal() *seqAttrGrpVal {
	return &seqAttrGrpVal{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrGrpVal) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrGrpVal) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Unique IDs for the grp_meta table
// create sequence seqAttrGrpMeta
// start 1024
type seqAttrGrpMeta struct {
	Val int
}

func newSeqAttrGrpMeta() *seqAttrGrpMeta {
	return &seqAttrGrpMeta{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqAttrGrpMeta) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqAttrGrpMeta) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//  Table: attrGrp
//
//  Description: A table to represent types of attributes.  A type is defined by
//               its subsystem, its grp ID and an attribute name.
// create table attrGrp
// constraint pk_attrGrp__id
//    pk(id)
type attrGrp struct {
	id      int    // integer
	subsys  string // varchar(256)
	name    string // varchar(256)
	sqlType string // varchar(30)
	active  bool
}

func newAttrGrp() *attrGrp {
	o := attrGrp{
		id:     0, // seqAttrGrp.NextVal(),
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrGrpVal
//
//  Description: A table to hold attribute values.
// create table attrGrpVal
// constraint pk_attrGrpVal__id
//    pk(id)
type attrGrpVal struct {
	id       int       // integer
	objectID int       // integer
	attrID   int       // integer
	dateVal  time.Time // timestamp
	shortVal string    // varchar(1024)
	blobVal  []byte    // text, like a blob
	serial   bool
	active   bool
}

func newAttrGrpVal() *attrGrpVal {
	o := attrGrpVal{
		id:     0, // seqAttrGrpVal.NextVal(),
		serial: false,
		active: true,
	}
	return &o
}

//  -------------------------------------------------------------------------------
//  Table: attrGrpMeta
//
//  Description: A table to represent metadata on types of attributes.
// create table attrGrpMeta
// constraint pk_attrGrpMeta__id
//    pk(id)
type attrGrpMeta struct {
	id     int    // integer
	attrID int    // integer
	name   string // varchar(256)
	value  string // varchar(2048)
	active bool
}

func newAttrGrpMeta() *attrGrpMeta {
	o := attrGrpMeta{
		id:     0, // seqAttrGrpMeta.NextVal(),
		active: true,
	}
	return &o
}

//  -----------------------------------------------------------------------------
//  Indexes.
//
//  Unique index on subsystem/name pair
//  Indexes on name and subsys.
//  Unique index on objectID/attrID pair
//  FK indexes on objectID and attrID.
//  Unique index on attrID/name pair
//  Index on meta name.
//  FK index on attrID.
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqJob
// start 1024
type seqJob struct {
	Val int
}

func newSeqJob() *seqJob {
	return &seqJob{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqJob) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqJob) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqJobMember
// start 1024
type seqJobMember struct {
	Val int
}

func newSeqJobMember() *seqJobMember {
	return &seqJobMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqJobMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqJobMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: job
//
// create table job
// constraint pk_jobID
//    pk(id)
type job struct {
	id              int       // integer
	name            []byte    // text, like a blob
	usrID           int       // integer
	schedTime       time.Time // timestamp
	priority        int       // integer
	compTime        time.Time // timestamp
	expire          bool
	failed          bool
	tries           int // integer
	executing       bool
	classID         int    // integer
	storyInstanceID int    // integer
	mediaInstanceID int    // integer
	errorMessage    []byte // text, like a blob
}

func newJob() *job {
	o := job{
		id:        0, // seqJob.NextVal(),
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

//
//  TABLE: jobResource
//
// create table jobResource
// constraint pk_jobResource
//    pk(jobID)
//    pk(resourceID)
type jobResource struct {
	jobID      int // integer
	resourceID int // integer
}

func newJobResource() *jobResource {
	o := jobResource{}
	return &o
}

//
//  TABLE: jobServerType
//
// create table jobServerType
// constraint pk_jobServerType
//    pk(jobID)
//    pk(serverTypeID)
type jobServerType struct {
	jobID        int // integer
	serverTypeID int // integer
}

func newJobServerType() *jobServerType {
	o := jobServerType{}
	return &o
}

//
//  TABLE: jobMember
//
// create table jobMember
// constraint pk_jobMemberID
//    pk(id)
type jobMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newJobMember() *jobMember {
	o := jobMember{
		id: 0, // seqJobMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: Michael Soderstrom <miraso@pacbell.net>
//
// 3487: /* Commented out because we're not using language stuff at this point.
// 3487:    By David.
// 3487:
// 3487:
// 3487: -- -----------------------------------------------------------------------------
// 3487: -- Sequences
// 3487:
// 3487: -- Unique IDs for the language table
// 3487: -- IDs under 1024 will contain dead languages
// 3487: CREATE SEQUENCE seq_language START 1024;
// 3487:
// 3487:
// 3487: -- -----------------------------------------------------------------------------
// 3487: -- Table: language
// 3487: --
// 3487: -- Description: name and description of languages
// 3487: --
// 3487:
// 3487: CREATE TABLE language (
// 3487:     id           INTEGER          NOT NULL
// 3487:                                 DEFAULT NEXTVAL('seq_language'),
// 3487:     name         VARCHAR(64),
// 3487:     description  VARCHAR(256),
// 3487:     active       BOOLEAN        NOT NULL DEFAULT TRUE,
// 3487:     CONSTRAINT pk_language__id PRIMARY KEY (id)
// 3487: );
// 3487:
// 3487: CREATE UNIQUE INDEX udx_language__name ON language(LOWER(name));
// 3487:
// 3487: */
//  Project: Bricolage
//
//  Target DBMS: PostgreSQL 7.1.2
//  Author: David Wheeler <david@justatheory.com>
//
//
//  SEQUENCES.
//
// create sequence seqMediaType
// start 1024
type seqMediaType struct {
	Val int
}

func newSeqMediaType() *seqMediaType {
	return &seqMediaType{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaType) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaType) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqMediaTypeExt
// start 1024
type seqMediaTypeExt struct {
	Val int
}

func newSeqMediaTypeExt() *seqMediaTypeExt {
	return &seqMediaTypeExt{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaTypeExt) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaTypeExt) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqMediaTypeMember
// start 1024
type seqMediaTypeMember struct {
	Val int
}

func newSeqMediaTypeMember() *seqMediaTypeMember {
	return &seqMediaTypeMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqMediaTypeMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqMediaTypeMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: mediaType
//
// create table mediaType
// constraint pk_mediaTypeID
//    pk(id)
type mediaType struct {
	id          int    // integer
	name        string // varchar(128)
	description string // varchar(256)
	active      bool
}

func newMediaType() *mediaType {
	o := mediaType{
		id:     0, // seqMediaType.NextVal(),
		active: true,
	}
	return &o
}

//
//  TABLE: mediaTypeExt
//
// create table mediaTypeExt
// constraint pk_mediaTypeExt__id
//    pk(id)
type mediaTypeExt struct {
	id          int    // integer
	mediaTypeID int    // integer
	extension   string // varchar(10)
}

func newMediaTypeExt() *mediaTypeExt {
	o := mediaTypeExt{
		id: 0, // seqMediaTypeExt.NextVal(),
	}
	return &o
}

//
//  TABLE: mediaTypeMember
//
// create table mediaTypeMember
// constraint pk_mediaTypeMember__id
//    pk(id)
type mediaTypeMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newMediaTypeMember() *mediaTypeMember {
	o := mediaTypeMember{
		id: 0, // seqMediaTypeMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//
//  Project: Bricolage API
//
//  Author: David Wheeler <david@justatheory.com>
//
//  SEQUENCES
//
// create sequence seqPref
// start 1024
type seqPref struct {
	Val int
}

func newSeqPref() *seqPref {
	return &seqPref{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPref) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPref) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqUsrPref
// start 1024
type seqUsrPref struct {
	Val int
}

func newSeqUsrPref() *seqUsrPref {
	return &seqUsrPref{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqUsrPref) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqUsrPref) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

// create sequence seqPrefMember
// start 1024
type seqPrefMember struct {
	Val int
}

func newSeqPrefMember() *seqPrefMember {
	return &seqPrefMember{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPrefMember) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPrefMember) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: pref
//         Global preferences.
// create table pref
// constraint pk_prefID
//    pk(id)
type pref struct {
	id              int    // integer
	name            string // varchar(64)
	description     string // varchar(256)
	value           string // varchar(256)
	def             string // varchar(256)
	manual          bool
	optType         string // varchar(16)
	canBeOverridden bool
}

func newPref() *pref {
	o := pref{
		id:              0, // seqPref.NextVal(),
		manual:          false,
		canBeOverridden: false,
	}
	return &o
}

//
//  TABLE: usrPref
//         Preferences overridden by a specific usr.
// create table usrPref
// constraint pk_usrPref__prefID__value
//    pk(id)
type usrPref struct {
	id     int    // integer
	prefID int    // integer
	usrID  int    // integer
	value  string // varchar(256)
}

func newUsrPref() *usrPref {
	o := usrPref{
		id: 0, // seqUsrPref.NextVal(),
	}
	return &o
}

//
//  TABLE: pref
//         Preference options.
// create table prefOpt
// constraint pk_prefOpt__prefID__value
//    pk(prefID)
//    pk(value)
type prefOpt struct {
	prefID      int    // integer
	value       string // varchar(256)
	description string // varchar(256)
}

func newPrefOpt() *prefOpt {
	o := prefOpt{}
	return &o
}

//
//  TABLE: prefMember
//
// create table prefMember
// constraint pk_prefMember__id
//    pk(id)
type prefMember struct {
	id       int // integer
	objectID int // integer
	memberID int // integer
}

func newPrefMember() *prefMember {
	o := prefMember{
		id: 0, // seqPrefMember.NextVal(),
	}
	return &o
}

//
//  INDEXES.
//
//
//  Project: Bricolage API
//
//  Author: David Wheeler <david@justatheory.com>
//
//  SEQUENCES
//
//  Use in both grpPriv and usr_priv (if we ever implement the latter).
// create sequence seqPriv
// start 1024
type seqPriv struct {
	Val int
}

func newSeqPriv() *seqPriv {
	return &seqPriv{Val: 1024}
}

// CurrVal returns the current value in the sequence.
func (seq *seqPriv) CurrVal() int {
	return seq.Val
}

// NextVal returns the next value in the sequence.
func (seq *seqPriv) NextVal() int {
	seq.Val = seq.Val + 1
	return seq.Val
}

//
//  TABLE: grpPriv
//         Privileges granted to user groups.
// create table grpPriv
// constraint pk_grpPrivID
//    pk(id)
type grpPriv struct {
	id    int       // integer
	grpID int       // integer
	value int       // integer
	mtime time.Time // timestamp
}

func newGrpPriv() *grpPriv {
	o := grpPriv{
		id:    0, // seqPriv.NextVal(),
		mtime: time.Now(),
	}
	// check value ( value  between  1  and  255 )
	return &o
}

//
//  TABLE: grpPrivGrpMember
//         Ties group privileges to groups for whose members the privilege
//         is granted.
// create table grpPrivGrpMember
// constraint pk_grpPrivGrpMember
//    pk(grpPrivID)
//    pk(grpID)
type grpPrivGrpMember struct {
	grpPrivID int // integer
	grpID     int // integer
}

func newGrpPrivGrpMember() *grpPrivGrpMember {
	o := grpPrivGrpMember{}
	return &o
}

//
//  INDEXES.
//
//  Everything below is left as a note for the future - in case we ever decided
//  actually allow privileges granted to individual users and/or individual
//  objects.
// 3719: /*
// 3719: --
// 3719: -- TABLE: grpPriv__grp
// 3719: --
// 3719:
// 3719: CREATE TABLE grpPriv__grp(
// 3719:     grpPrivID    INTEGER           NOT NULL,
// 3719:     grpID         INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_grpPriv__grp PRIMARY KEY (grpPrivID,grpID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: grpPriv__person
// 3719: --
// 3719:
// 3719: CREATE TABLE grpPriv__person(
// 3719:     grpPrivID    INTEGER           NOT NULL,
// 3719:     personID      INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_grpPriv__person PRIMARY KEY (grpPrivID,personID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: grpPriv__usr
// 3719: --
// 3719:
// 3719: CREATE TABLE grpPriv__usr(
// 3719:     grpPrivID    INTEGER           NOT NULL,
// 3719:     usrID        INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_grpPriv__usr PRIMARY KEY (grpPrivID,usrID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: priv_table
// 3719: --
// 3719:
// 3719: CREATE TABLE priv_table(
// 3719:     id      INTEGER           NOT NULL,
// 3719:     name    VARCHAR(30)    NOT NULL,
// 3719:     CONSTRAINT pk_priv_table__id PRIMARY KEY (id)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: usr_priv
// 3719: --
// 3719:
// 3719: CREATE TABLE usr_priv(
// 3719:     id          INTEGER           NOT NULL,
// 3719:     usrID    INTEGER           NOT NULL,
// 3719:     value       INT2     NOT NULL,
// 3719:     CONSTRAINT pk_usr_priv__id PRIMARY KEY (id)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: usr_priv__grp
// 3719: --
// 3719:
// 3719: CREATE TABLE usr_priv__grp(
// 3719:     priv_usrID    INTEGER           NOT NULL,
// 3719:     grpID          INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_usr_priv__grp PRIMARY KEY (priv_usrID,grpID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: usr_priv__person
// 3719: --
// 3719:
// 3719: CREATE TABLE usr_priv__person(
// 3719:     usr_priv__id    INTEGER           NOT NULL,
// 3719:     personID       INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_usr_priv__person PRIMARY KEY (usr_priv__id,personID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: usr_priv__usr
// 3719: --
// 3719:
// 3719: CREATE TABLE usr_priv__usr(
// 3719:     usr_priv__id    INTEGER           NOT NULL,
// 3719:     usrID         INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_usr_priv__usr PRIMARY KEY (usr_priv__id,usrID)
// 3719: )
// 3719: ;
// 3719:
// 3719:
// 3719: --
// 3719: -- TABLE: usr_priv__grpMember
// 3719: --
// 3719:
// 3719: CREATE TABLE usr_priv__grpMember(
// 3719:     usr_priv__id    INTEGER           NOT NULL,
// 3719:     grpID          INTEGER           NOT NULL,
// 3719:     CONSTRAINT pk_usr_priv__grpMember PRIMARY KEY (usr_priv__id,grpID)
// 3719: )
// 3719: ;
// 3719:
// 3719: */
