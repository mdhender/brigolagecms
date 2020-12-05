-- Project: Bricolage
--
-- Author: David Wheeler <david@justatheory.com>
--

-- This DDL is for the creation of universal stuff needed by other DDLs, such as
-- functions.

--
-- Functions. 
--

-- This function allows us to create UNIQUE indices that combine a lowercased
-- TEXT (or VARCHAR) column with an INTEGER column. See Bric/Util/AlertType.sql
-- for an example.
CREATE   FUNCTION lower_text_num(TEXT, INTEGER)
RETURNS  TEXT AS 'SELECT LOWER($1) || to_char($2, ''|FM9999999999'')'
LANGUAGE 'sql'
WITH     (ISCACHABLE);

/* XXX Once we require 7.4 or later (2.0?), dump the function and aggregate
 * below in favor of this aggregate:

CREATE AGGREGATE array_accum (
    sfunc    = array_append,
    basetype = anyelement,
    stype    = anyarray,
    initcond = '{}'
);

 * Then change the usage from group_concat(foo) to
 * array_to_string(array_accum(foo), ' ')

*/

-- This function is used to append a space followed by a number to a TEXT
-- string. It is used primarily for the group_concat aggregate (below). We omit
-- the ID 0 because it is a hidden, secret group to which permissions do not
-- apply.
CREATE   FUNCTION append_id(TEXT, INTEGER)
RETURNS  TEXT AS '
    SELECT CASE WHEN $2 = 0 THEN
                $1
           ELSE
                $1 || '' '' || CAST($2 AS TEXT)
           END;'
LANGUAGE 'sql'
WITH     (ISCACHABLE, ISSTRICT);

-- This aggregate is designed to concatenate all of the IDs that would
-- otherwise cause a query to return multiple rows into a single value
-- with each ID separated by a space. This makes it easy for us to pull
-- out the list of IDs using split, _and_ keeps the entire contents of
-- an object in a single row, thus also enabling the use of OFFSET and
-- LIMIT.
CREATE AGGREGATE group_concat (
    SFUNC    = append_id,
    BASETYPE = INTEGER,
    STYPE    = TEXT,
    INITCOND = ''
);

/*
-- This is a temporary compatibility measure.
CREATE FUNCTION int_to_boolean(integer) RETURNS boolean
  AS 'select case when $1 = 0 then false else true end'
LANGUAGE 'sql' IMMUTABLE;

CREATE CAST (integer AS boolean)
  WITH FUNCTION int_to_boolean(integer) AS IMPLICIT;
*/-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- This is the SQL that will create the element table.
-- It is related to the Bric::ElementType class.
--
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the element table
CREATE SEQUENCE seq_at_type START 1024;
CREATE SEQUENCE seq_at_type_member START 1024;

-- -----------------------------------------------------------------------------
-- Table: element
--
-- Description:    The table that holds the information for a given asset type.  
--         Holds name and description information and is references by 
--        element_type rows.
--

CREATE TABLE at_type (
    id              INTEGER        NOT NULL
                                   DEFAULT NEXTVAL('seq_at_type'),
    name            VARCHAR(64)       NOT NULL,
    description     VARCHAR(256),
    top_level       BOOLEAN        NOT NULL DEFAULT FALSE,
    paginated       BOOLEAN        NOT NULL DEFAULT FALSE,
    fixed_url       BOOLEAN        NOT NULL DEFAULT FALSE,
    related_story   BOOLEAN        NOT NULL DEFAULT FALSE,
    related_media   BOOLEAN        NOT NULL DEFAULT FALSE,
    media           BOOLEAN        NOT NULL DEFAULT FALSE,
    biz_class__id   INTEGER        NOT NULL,
    active          BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_at_type__id PRIMARY KEY (id)
);

--
-- TABLE: element_type_member
--

CREATE TABLE at_type_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_at_type_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_at_type_member__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Indexes.
--
CREATE UNIQUE INDEX udx_at_type__name ON at_type(LOWER(name));
CREATE INDEX fkx_class__at_type ON at_type(biz_class__id);
CREATE INDEX fkx_at_type__at_type_member ON at_type_member(object_id);
CREATE INDEX fkx_member__at_type_member ON at_type_member(member__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- The sql representation of Media Assets.

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the media table
CREATE SEQUENCE seq_media START 1024;

CREATE SEQUENCE seq_media_instance START 1024;

-- Unique ids for the media__contributor table
CREATE SEQUENCE seq_media__contributor START 1024;

-- Unique IDs for the media_member table
CREATE SEQUENCE seq_media_member START 1024;

CREATE SEQUENCE seq_media_fields START 1024;

-- Unique IDs for the media_uri table
CREATE SEQUENCE seq_media_uri START 1024;


-- -----------------------------------------------------------------------------
-- Table media
-- 
-- Description: The Media table this houses the data for a given media asset
--                              and its related asset_version_data
--

CREATE TABLE media (
    id                INTEGER   NOT NULL
                                      DEFAULT NEXTVAL('seq_media'),
    uuid              TEXT            NOT NULL,
    element_type__id  INTEGER   NOT NULL,
    current_version   INTEGER,
    published_version INTEGER,
    usr__id           INTEGER,
    first_publish_date TIMESTAMP,
    publish_date      TIMESTAMP,
    workflow__id      INTEGER   NOT NULL,
    desk__id          INTEGER   NOT NULL,
    publish_status    BOOLEAN    NOT NULL DEFAULT FALSE,
    active            BOOLEAN    NOT NULL DEFAULT TRUE,
    site__id          INTEGER   NOT NULL,
    alias_id          INTEGER   CONSTRAINT ck_media_id
                                        CHECK (alias_id != id),  
    CONSTRAINT pk_media__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: media_instance
--
-- Description: An instance of a media object
--
--

CREATE TABLE media_instance (
    id                  INTEGER   NOT NULL
                                        DEFAULT NEXTVAL('seq_media_instance'),
    name                VARCHAR(256),
    description         VARCHAR(1024),
    media__id           INTEGER   NOT NULL,
    source__id          INTEGER   NOT NULL,
    priority            INT2      NOT NULL
                                  DEFAULT 3
                                  CONSTRAINT ck_media__priority
                                  CHECK (priority BETWEEN 1 AND 5),
    usr__id             INTEGER   NOT NULL,
    version             INTEGER,
    expire_date         TIMESTAMP,
    category__id        INTEGER   NOT NULL,
    media_type__id      INTEGER   NOT NULL,
    primary_oc__id      INTEGER   NOT NULL,
    file_size           INTEGER,
    file_name           VARCHAR(256),
    location            VARCHAR(256),
    uri                 VARCHAR(256),
    cover_date          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note                TEXT,
    checked_out         BOOLEAN    NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_media_instance__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table media_uri
--
-- Description: Tracks all URIs for stories.
--
CREATE TABLE media_uri (
    id        INTEGER     NOT NULL
                              DEFAULT NEXTVAL('seq_media_uri'),
    media__id INTEGER     NOT NULL,
    site__id  INTEGER     NOT NULL,
    uri       TEXT            NOT NULL,
    CONSTRAINT pk_media_uri__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table media__output_channel
-- 
-- Description: Mapping Table between stories and output channels.
--
--

CREATE TABLE media__output_channel (
    media_instance__id  INTEGER  NOT NULL,
    output_channel__id  INTEGER  NOT NULL,
    CONSTRAINT pk_media_output_channel
      PRIMARY KEY (media_instance__id, output_channel__id)
);

-- -----------------------------------------------------------------------------
-- Table: media_fields
-- 
-- Description: A mapping table between Media classes and functions that
--                              Will be run against uploaded files
-- 
CREATE TABLE media_fields (
    id              INTEGER  NOT NULL     
                                   DEFAULT NEXTVAL('seq_media_fields'),
    biz_pkg         INTEGER  NOT NULL,
    name            VARCHAR(32)    NOT NULL,
    function_name   VARCHAR(256)   NOT NULL,
    active          BOOLEAN   NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_media_fields__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table media__contributor
-- 
-- Description: mapping tables between media instances and contributors
--
--

CREATE TABLE media__contributor (
    id                  INTEGER   NOT NULL
                                        DEFAULT NEXTVAL('seq_media__contributor'),
    media_instance__id  INTEGER   NOT NULL,
    member__id          INTEGER   NOT NULL,
    place               INT2      NOT NULL,
    role                VARCHAR(256),
    CONSTRAINT pk_media_category_id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: media_member
-- 
-- Description: The link between media objects and member objects
--

CREATE TABLE media_member (
    id          INTEGER  NOT NULL
                               DEFAULT NEXTVAL('seq_media_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_media_member__id PRIMARY KEY (id)
);


--
-- INDEXES.
--

-- media
CREATE INDEX idx_media__uuid ON media(uuid);
CREATE INDEX idx_media__first_publish_date ON media(first_publish_date);
CREATE INDEX idx_media__publish_date ON media(publish_date);
CREATE INDEX idx_media_instance__cover_date ON media_instance(cover_date);
CREATE INDEX fkx_usr__media ON media(usr__id);
CREATE INDEX fkx_element_type__media ON media(element_type__id);
CREATE INDEX fkx_site_id__media ON media(site__id);
CREATE INDEX fkx_alias_id__media ON media(alias_id);

-- media_instance
CREATE INDEX idx_media_instance__name ON media_instance(LOWER(name));
CREATE INDEX idx_media_instance__description ON media_instance(LOWER(description));
CREATE INDEX fkx_media_instance__source ON media_instance(source__id);
CREATE INDEX idx_media_instance__file_name ON media_instance(LOWER(file_name));
CREATE INDEX idx_media_instance__uri ON media_instance(LOWER(uri));
CREATE UNIQUE INDEX udx_media__media_instance ON media_instance(media__id, version, checked_out);
CREATE INDEX fkx_usr__media_instance ON media_instance(usr__id);
CREATE INDEX fkx_media_type__media_instance ON media_instance(media_type__id);
CREATE INDEX fkx_category__media_instance ON media_instance(category__id);
CREATE INDEX fkx_primary_oc__media_instance ON media_instance(primary_oc__id);
CREATE INDEX idx_media_instance__note ON media_instance(note) WHERE note IS NOT NULL;

-- media_uri
CREATE INDEX fkx_media__media_uri ON media_uri(media__id);
CREATE UNIQUE INDEX udx_media_uri__site_id__uri
ON media_uri(lower_text_num(uri, site__id));

-- media__output_channel
CREATE INDEX fkx_media__oc__media ON media__output_channel(media_instance__id);
CREATE INDEX fkx_media__oc__oc ON media__output_channel(output_channel__id);

--media__contributor
CREATE INDEX fkx_media__media__contributor ON media__contributor(media_instance__id);
CREATE INDEX fkx_member__media__contributor ON media__contributor(member__id);

-- media_member.
CREATE INDEX fkx_media__media_member ON media_member(object_id);
CREATE INDEX fkx_member__media_member ON media_member(member__id);

CREATE INDEX fkx_media__desk__id ON media(desk__id) WHERE desk__id > 0;
CREATE INDEX fkx_media__workflow__id ON media(workflow__id) WHERE workflow__id > 0;
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
--
-- This is the SQL representation of the story object
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the story table
CREATE SEQUENCE seq_story START  1024;

-- Unique IDs for the story_instance table
CREATE SEQUENCE seq_story_instance START 1024;

-- Unique IDs for the story__category mapping table
CREATE SEQUENCE seq_story__category START  1024;

-- Unique ids for the story__contributor table
CREATE SEQUENCE seq_story__contributor START 1024;

-- Unique IDs for the story_uri table
CREATE SEQUENCE seq_story_uri START 1024;


-- -----------------------------------------------------------------------------
-- Table: story
--
-- Description: The story properties. Versioning info might get added here and
--              the rights info might get removed. It is also possible that
--              the asset type field will need a cascading delete.


CREATE TABLE story (
    id                INTEGER         NOT NULL
                                      DEFAULT NEXTVAL('seq_story'),
    uuid              TEXT            NOT NULL,
    usr__id           INTEGER,
    element_type__id  INTEGER         NOT NULL,
    first_publish_date TIMESTAMP,
    publish_date      TIMESTAMP,
    current_version   INTEGER         NOT NULL,
    published_version INTEGER,
    workflow__id      INTEGER         NOT NULL,
    desk__id          INTEGER         NOT NULL,
    publish_status    BOOLEAN         NOT NULL DEFAULT FALSE,
    active            BOOLEAN         NOT NULL DEFAULT TRUE,
    site__id          INTEGER         NOT NULL,
    alias_id          INTEGER         CONSTRAINT ck_story_id
                                        CHECK (alias_id != id),  
    CONSTRAINT pk_story__id PRIMARY KEY (id)
);

-- ----------------------------------------------------------------------------
-- Table story_instance
--
-- Description:  An instance of a story
--

CREATE TABLE story_instance (
    id             INTEGER      NOT NULL
                                DEFAULT NEXTVAL('seq_story_instance'),
    name           VARCHAR(256),
    description    VARCHAR(1024),
    story__id      INTEGER      NOT NULL,
    source__id     INTEGER      NOT NULL,
    primary_uri    VARCHAR(128),
    priority       INT2         NOT NULL
                                DEFAULT 3
                                CONSTRAINT ck_story__priority
                                CHECK (priority BETWEEN 1 AND 5),
    version        INTEGER,
    expire_date    TIMESTAMP,
    usr__id        INTEGER      NOT NULL,
    slug           VARCHAR(64),
    primary_oc__id INTEGER      NOT NULL,
    cover_date     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note           TEXT,
    checked_out    BOOLEAN      NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_story_instance__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table story_uri
--
-- Description: Tracks all URIs for stories.
--
CREATE TABLE story_uri (
    id        INTEGER     NOT NULL
                              DEFAULT NEXTVAL('seq_story_uri'),
    story__id INTEGER     NOT NULL,
    site__id INTEGER      NOT NULL,
    uri       TEXT            NOT NULL,
    CONSTRAINT pk_story_uri__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table story__output_channel
-- 
-- Description: Mapping Table between stories and output channels.
--
--

CREATE TABLE story__output_channel (
    story_instance__id  INTEGER  NOT NULL,
    output_channel__id  INTEGER  NOT NULL,
    CONSTRAINT pk_story_output_channel
      PRIMARY KEY (story_instance__id, output_channel__id)
);


-- -----------------------------------------------------------------------------
-- Table story__category
-- 
-- Description: Mapping Table between Stories and categories
--
--

CREATE TABLE story__category (
    id                  INTEGER  NOT NULL
                                       DEFAULT NEXTVAL('seq_story__category'),
    story_instance__id  INTEGER  NOT NULL,
    category__id        INTEGER  NOT NULL,
    main                BOOLEAN   NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_story_category__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table story__contributor
-- 
-- Description: mapping tables between story instances and contributors
--
--

CREATE TABLE story__contributor (
    id                  INTEGER   NOT NULL
                                        DEFAULT NEXTVAL('seq_story__contributor'),
    story_instance__id  INTEGER   NOT NULL,
    member__id          INTEGER   NOT NULL,
    place               INT2      NOT NULL,
    role                VARCHAR(256),
    CONSTRAINT pk_story_category_id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--

-- story
CREATE INDEX idx_story__uuid ON story(uuid);
CREATE INDEX fkx_usr__story ON story(usr__id);
CREATE INDEX fkx_element_type__story ON story(element_type__id);
CREATE INDEX fkx_site_id__story ON story(site__id);
CREATE INDEX fkx_alias_id__story ON story(alias_id);
CREATE INDEX idx_story__first_publish_date ON story(first_publish_date);
CREATE INDEX idx_story__publish_date ON story(publish_date);
CREATE INDEX idx_story_instance__cover_date ON story_instance(cover_date);

-- story_instance
CREATE INDEX idx_story_instance__name ON story_instance(LOWER(name));
CREATE INDEX idx_story_instance__description ON story_instance(LOWER(description));
CREATE INDEX fkx_story_instance__source__id ON story_instance(source__id);
CREATE INDEX idx_story_instance__slug ON story_instance(LOWER(slug));
CREATE INDEX idx_story_instance__primary_uri ON story_instance(LOWER(primary_uri));
CREATE UNIQUE INDEX udx_story__story_instance ON story_instance(story__id, version, checked_out);
CREATE INDEX fkx_usr__story_instance ON story_instance(usr__id);
CREATE INDEX fkx_primary_oc__story_instance ON story_instance(primary_oc__id);
CREATE INDEX idx_story_instance__note ON story_instance(note) WHERE note IS NOT NULL;

-- story_uri
CREATE INDEX fkx_story__story_uri ON story_uri(story__id);
CREATE UNIQUE INDEX udx_story_uri__site_id__uri
ON story_uri(lower_text_num(uri, site__id));

-- story__category
CREATE UNIQUE INDEX udx_story_category__story__cat ON story__category(story_instance__id, category__id);
CREATE INDEX fkx_story__story__category ON story__category(story_instance__id);
CREATE INDEX fkx_category__story__category ON story__category(category__id);

-- story__output_channel
CREATE INDEX fkx_story__oc__story ON story__output_channel(story_instance__id);
CREATE INDEX fkx_story__oc__oc ON story__output_channel(output_channel__id);

--story__contributor
CREATE INDEX fkx_story__story__contributor ON story__contributor(story_instance__id);
CREATE INDEX fkx_member__story__contributor ON story__contributor(member__id);

CREATE INDEX fkx_story__desk__id ON story(desk__id) WHERE desk__id > 0;
CREATE INDEX fkx_story__workflow__id ON story(workflow__id) WHERE workflow__id > 0;
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- The sql that will hold all the template asset information.

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique ids for the template table
CREATE SEQUENCE seq_template START 1024;

CREATE SEQUENCE seq_template_instance START 1024;

-- Unique IDs for the story_member table
CREATE SEQUENCE seq_template_member START 1024;


-- -----------------------------------------------------------------------------
-- Table template
--
-- Description: The table that holds all the template info
--
--
CREATE TABLE template (
    id                  INTEGER        NOT NULL
                                       DEFAULT NEXTVAL('seq_template'),
    usr__id             INTEGER,  
    output_channel__id  INTEGER        NOT NULL,
    tplate_type         INT2           NOT NULL
                                       DEFAULT 1
                                       CONSTRAINT ck_template___tplate_type
                                         CHECK (tplate_type IN (1, 2, 3)),
    element_type__id    INTEGER,
    file_name           TEXT,
    current_version     INTEGER        NOT NULL,
    workflow__id        INTEGER        NOT NULL,
    desk__id            INTEGER        NOT NULL,
    published_version   INTEGER,
    deploy_status       BOOLEAN        NOT NULL DEFAULT FALSE,
    deploy_date         TIMESTAMP,
    active              BOOLEAN        NOT NULL DEFAULT TRUE,
    site__id            INTEGER        NOT NULL,
    CONSTRAINT pk_template__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table template_instance
--
-- Description:  An versioned instance of a template asset
--

CREATE TABLE template_instance (
    id              INTEGER        NOT NULL
                                   DEFAULT NEXTVAL('seq_template_instance'),
    name            VARCHAR(256),
    description     VARCHAR(1024),
    priority        INT2           NOT NULL
                                   DEFAULT 3
                                   CONSTRAINT ck_story__priority
                                   CHECK (priority BETWEEN 1 AND 5),
    template__id    INTEGER        NOT NULL,
    category__id    INTEGER,
    version         INTEGER,
    usr__id         INTEGER        NOT NULL,
    file_name       TEXT,
    data            TEXT,
    expire_date     TIMESTAMP,
    note            TEXT,
    checked_out     BOOLEAN        NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_template_instance__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: template_member
-- 
-- Description: The link between template objects and member objects
--

CREATE TABLE template_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_template_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_template_member__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--

-- template.
CREATE UNIQUE INDEX udx_template__file_name__oc
       ON template(file_name, output_channel__id);
CREATE INDEX idx_template__file_name ON template(LOWER(file_name));
CREATE INDEX fkx_usr__template ON template(usr__id);
CREATE INDEX fkx_output_channel__template ON template(output_channel__id);
CREATE INDEX fkx_element_type__template ON template(element_type__id);
CREATE INDEX fkx_template__desk__id ON template(desk__id) WHERE desk__id > 0;
CREATE INDEX fkx_template__workflow__id ON template(workflow__id) WHERE workflow__id > 0;
CREATE INDEX fkx_site__template ON template(site__id);

-- template_instance.
CREATE INDEX idx_template_instance__name ON template_instance(LOWER(name));
CREATE INDEX fkx_usr__template_instance ON template_instance(usr__id);
CREATE INDEX fkx_template__tmpl_instance ON template_instance(template__id);
CREATE INDEX idx_template_instance__note ON template_instance(note) WHERE note IS NOT NULL;
CREATE INDEX fkx_template_instance_category ON template_instance(category__id);

-- template_member.
CREATE INDEX fkx_template__template_member ON template_member(object_id);
CREATE INDEX fkx_member__template_member ON template_member(member__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- This SQL creates the tables necessary for the attribute object.  This file
-- applies to attributes on the Bric::Person class.  Any other classes that 
-- require attributes need only duplicate these tables, changing 'person' to 
-- the correct class name.  Class names may be shortened to ensure that the
-- resulting table names are under the oracle 30 character name limit as long
-- as the resulting shortened class name is unique.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the category table
CREATE SEQUENCE seq_category START  1024;

-- Unique IDs for the category_member table
CREATE SEQUENCE seq_category_member START  1024;

-- Unique IDs for the attr_category table
CREATE SEQUENCE seq_attr_category START 1024;

-- Unique IDs for each attr_category_val table
CREATE SEQUENCE seq_attr_category_val START 1024;

-- Unique IDs for the category_meta table
CREATE SEQUENCE seq_attr_category_meta START 1024;

-- -----------------------------------------------------------------------------
-- Table: category
-- 
-- Description: The category table
--

CREATE TABLE category (
    id               INTEGER         NOT NULL
                                     DEFAULT NEXTVAL('seq_category'),
    site__id         INTEGER         NOT NULL,
    directory        VARCHAR(128)    NOT NULL,
    uri              VARCHAR(256)    NOT NULL,
    name             VARCHAR(128),
    description      VARCHAR(256),
    parent_id        INTEGER         NOT NULL,
    asset_grp_id     INTEGER         NOT NULL,
    active           BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_category__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: category_member
-- 
-- Description: The link between desk objects and member objects
--

CREATE TABLE category_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_category_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_category_member__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_category
--
-- Description: A table to represent types of attributes.  A type is defined by
--              its subsystem, its category ID and an attribute name.

CREATE TABLE attr_category (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_category'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_category__id PRIMARY KEY (id)
);



-- -----------------------------------------------------------------------------
-- Table: attr_category_val
--
-- Description: A table to hold attribute values.

CREATE TABLE attr_category_val (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_attr_category_val'),
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_category_val__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: attr_category_meta
--
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_category_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_category_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_category_meta__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Indexes.
--
CREATE INDEX idx_category__directory ON category(LOWER(directory));
CREATE UNIQUE INDEX udx_category__site_uri ON category(site__id, uri);
CREATE INDEX idx_category__uri ON category(uri);
CREATE INDEX idx_category__lower_uri ON category(LOWER(uri));
CREATE INDEX idx_category__name ON category(LOWER(name));
CREATE INDEX idx_category__parent_id ON category(parent_id);
CREATE INDEX fkx_asset_grp__category ON category(asset_grp_id);

CREATE INDEX fkx_category__category_member ON category_member(object_id);
CREATE INDEX fkx_member__category_member ON category_member(member__id);
CREATE INDEX fkx_category__site ON category(site__id);

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_cat__subsys__name ON attr_category(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_cat__name ON attr_category(LOWER(name));
CREATE INDEX idx_attr_cat__subsys ON attr_category(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_cat_val__obj_attr ON attr_category_val (object__id,attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_cat__attr_cat_val ON attr_category_val(object__id);
CREATE INDEX fkx_attr_cat__attr_cat_val ON attr_category_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_cat_meta__attr_name ON attr_category_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_cat_meta__name ON attr_category_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_cat__attr_cat_meta ON attr_category_meta(attr__id);

-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- This DDL creates the basic tables for all Bric::BC::Contact objects.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_contact START 1024;
CREATE SEQUENCE seq_contact_value START 1024;

-- 
-- TABLE: contact 
--

CREATE TABLE contact (
    id           INTEGER           NOT NULL
                                   DEFAULT NEXTVAL('seq_contact'),
    type         VARCHAR(64)       NOT NULL,
    description     VARCHAR(256),
    active       BOOLEAN           NOT NULL DEFAULT TRUE,
    alertable    BOOLEAN           NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_contact__id PRIMARY KEY (id)
);

-- 
-- TABLE: contact_value
--

CREATE TABLE contact_value (
    id           INTEGER           NOT NULL
                                   DEFAULT NEXTVAL('seq_contact_value'),
    contact__id  INTEGER           NOT NULL,
    value         VARCHAR(256)       NOT NULL,
    active       BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_contact_value__id PRIMARY KEY (id)
);

--
-- TABLE: person__contact
--
CREATE TABLE person__contact_value (
    person__id          INTEGER    NOT NULL,
    contact_value__id   INTEGER    NOT NULL,
    CONSTRAINT pk_person__contact_value PRIMARY KEY (person__id, contact_value__id)
);


-- 
-- INDEXES.
--

CREATE UNIQUE INDEX udx_contact__type ON contact(LOWER(type));
CREATE INDEX idx_contact_value_value ON contact_value(LOWER(value));
CREATE INDEX fkx_contact__contact_value on contact_value(contact__id);

CREATE INDEX fkx_person__p_c_val ON person__contact_value(person__id);
CREATE INDEX fkx_contact_value__p_c_val ON person__contact_value(contact_value__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- This is the sql that will create the container elements
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the story element table
CREATE SEQUENCE seq_story_element START  1024;

-- Unique IDs for the media element table
CREATE SEQUENCE seq_media_element START  1024;

-- -----------------------------------------------------------------------------
-- Table story_element
-- 
-- Description: Holds the properties of a container element. Note that
--              elements can hold either other container elements or field
--              elements, not both.
--

CREATE TABLE story_element (
    id                   INTEGER         NOT NULL
                                         DEFAULT NEXTVAL('seq_story_element'),
    element_type__id     INTEGER         NOT NULL,
    object_instance_id   INTEGER         NOT NULL,
    parent_id            INTEGER,
    place                INTEGER         NOT NULL,
    object_order         INTEGER         NOT NULL,
    displayed            BOOLEAN         NOT NULL DEFAULT FALSE,
    related_story__id    INTEGER,
    related_media__id    INTEGER,
    active               BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_story_element__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table media_element
--
-- Description: Holds the properties of a media container element.
--
--

CREATE TABLE media_element (
    id                          INTEGER         NOT NULL
                                                DEFAULT NEXTVAL('seq_media_element'),
    element_type__id            INTEGER         NOT NULL,
    object_instance_id          INTEGER         NOT NULL,
    parent_id                   INTEGER,
    place                       INTEGER         NOT NULL,
    object_order                INTEGER         NOT NULL,
    displayed                   BOOLEAN         NOT NULL DEFAULT FALSE,
    related_story__id           INTEGER, 
    related_media__id           INTEGER,
    active                      BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_media_element__id PRIMARY KEY (id)
);


--
-- INDEXES.
--
CREATE INDEX fkx_story_element__story_element ON story_element(parent_id);
CREATE INDEX fkx_story__story_element ON story_element(object_instance_id);
CREATE INDEX fkx_story_element__related_story ON story_element(related_story__id);
CREATE INDEX fkx_story_element__related_media ON story_element(related_media__id);
CREATE INDEX fkx_story_element__element_type ON story_element(element_type__id);

CREATE INDEX fkx_media_element__media_element ON media_element(parent_id);
CREATE INDEX fkx_media__media_element ON media_element(object_instance_id);
CREATE INDEX fkx_media_element__related_story ON media_element(related_story__id);
CREATE INDEX fkx_media_element__related_media ON media_element(related_media__id);
CREATE INDEX fkx_media_element__element_type ON media_element(element_type__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the story field elements.
CREATE SEQUENCE seq_story_field START 1024;

-- Unique IDs for the media field elements.
CREATE SEQUENCE seq_media_field START 1024;

-- -----------------------------------------------------------------------------
-- Table story_field
--
-- Description: Story Field elements are story specific mappings to the
--              Bric::Biz::Element::Field class. They link to the story that
--              this element is a part of, the attribute id of the data that
--              is contained with in, and it's parent's id (a story_element
--              row). Place is it's order and active is it's active state.
--
--

CREATE TABLE story_field (
    id                   INTEGER        NOT NULL
                                        DEFAULT NEXTVAL('seq_story_field'),
    field_type__id       INTEGER        NOT NULL,
    object_instance_id   INTEGER        NOT NULL,
    parent_id            INTEGER        NOT NULL,
    hold_val             BOOLEAN        NOT NULL DEFAULT FALSE,
    place                INTEGER        NOT NULL,
    object_order         INTEGER        NOT NULL,
    date_val             TIMESTAMP,
    short_val            TEXT,
    blob_val             TEXT,
    active               BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_story_field__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table media_field
--
-- Description: Media Field elements are media specific mappings to the
--              Bric::Biz::Element::Field class. They link to the media that
--              this element is a part of, the attribute id of the data that
--              is contained with in, and it's parent's id (a media_element
--              row). Place is it's order and active is it's active state.
--
--

CREATE TABLE media_field (
    id                   INTEGER        NOT NULL
                                        DEFAULT NEXTVAL('seq_media_field'),
    field_type__id       INTEGER        NOT NULL,
    object_instance_id   INTEGER        NOT NULL,
    parent_id            INTEGER        NOT NULL,
    place                INTEGER        NOT NULL,
    hold_val             BOOLEAN        NOT NULL DEFAULT FALSE,
    object_order         INTEGER        NOT NULL,
    date_val             TIMESTAMP,
    short_val            VARCHAR(1024),
    blob_val             TEXT,
    active               BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_media_field__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_story_instance__story_field ON story_field(object_instance_id);
CREATE INDEX fkx_field_type__story_field ON story_field(field_type__id);
CREATE INDEX fkx_story_field__story_field ON story_field(parent_id);

CREATE INDEX fkx_media_instance__media_field ON media_field(object_instance_id);
CREATE INDEX fkx_field_type__media_field ON media_field(field_type__id);
CREATE INDEX fkx_media_field__media_field ON media_field(parent_id);
-- Project: Bricolage
--
--
-- This is the SQL that will create the element_type table.
-- It is related to the Bric::ElementType class.
-- Related tables are element and field
--
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the element table
CREATE SEQUENCE seq_element_type START 1024;

-- Unique IDs for the subelement table
CREATE SEQUENCE seq_subelement_type START 1024;

-- Unique IDs for element__output_channel
CREATE SEQUENCE seq_element_type__output_channel START 1024;

-- Unique IDs for element__language
CREATE SEQUENCE seq_element__language START 1024;

-- Unique IDs for element_type_member
CREATE SEQUENCE seq_element_type_member START 1024;

-- Unique IDs for the attr_element table
CREATE SEQUENCE seq_attr_element_type START 1024;

-- Unique IDs for each attr_element_*_val table
CREATE SEQUENCE seq_attr_element_type_val START 1024;

-- Unique IDs for the element_meta table
CREATE SEQUENCE seq_attr_element_type_meta START 1024;

-- Unique IDs for the element__site table.
CREATE SEQUENCE seq_element_type__site START 1024;

-- -----------------------------------------------------------------------------
-- Table: element_type
--
-- Description: The table that holds the information for a given asset type.  
--              Holds name and description information and is references by 
--              element_contaner and field_type rows.
--

CREATE TABLE element_type  (
    id              INTEGER        NOT NULL
                                   DEFAULT NEXTVAL('seq_element_type'),
    name            VARCHAR(64)    NOT NULL,
    key_name        VARCHAR(64)    NOT NULL,
    description     VARCHAR(256),
    top_level       BOOLEAN        NOT NULL DEFAULT FALSE,
    paginated       BOOLEAN        NOT NULL DEFAULT FALSE,
    fixed_uri       BOOLEAN        NOT NULL DEFAULT FALSE,
    related_story   BOOLEAN        NOT NULL DEFAULT FALSE,
    related_media   BOOLEAN        NOT NULL DEFAULT FALSE,
    displayed       BOOLEAN        NOT NULL DEFAULT FALSE,
    media           BOOLEAN        NOT NULL DEFAULT FALSE,
    biz_class__id   INTEGER        NOT NULL,
    type__id        INTEGER,
    active          BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_element_type__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: subelement_type
--
-- Description: A table that manages element type parent/child relationships.

CREATE TABLE subelement_type  (
    id              INTEGER        NOT NULL
                                   DEFAULT NEXTVAL('seq_subelement_type'),
    parent_id       INTEGER        NOT NULL,
    child_id        INTEGER        NOT NULL,
    place           INTEGER        NOT NULL DEFAULT 1,
    min_occurrence  INTEGER        NOT NULL DEFAULT 0,
    max_occurrence  INTEGER        NOT NULL DEFAULT 0,
    CONSTRAINT pk_subelement_type__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: element__site
--
-- Description: A table that maps 

CREATE TABLE element_type__site (
    id               INTEGER NOT NULL
                             DEFAULT NEXTVAL('seq_element_type__site'),
    element_type__id INTEGER NOT NULL,
    site__id         INTEGER NOT NULL,
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    primary_oc__id   INTEGER NOT NULL,
    CONSTRAINT pk_element_type__site__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: element__output_channel
--
-- Description: Holds a reference to the asset type table, the output channel 
--              table and an active flag
--

CREATE TABLE element_type__output_channel (
    id                  INTEGER    NOT NULL
                                   DEFAULT NEXTVAL('seq_element_type__output_channel'),
    element_type__id    INTEGER    NOT NULL,
    output_channel__id  INTEGER    NOT NULL,
    enabled             BOOLEAN    NOT NULL DEFAULT TRUE,
    active              BOOLEAN    NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_element_type__output_channel__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: element_type_member
-- 
-- Description: The link between element objects and member objects
--
CREATE TABLE element_type_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_element_type_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_element_type_member__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_element
--
-- Description: A table to represent types of attributes.  A type is defined by
--              its subsystem, its element ID and an attribute name.

CREATE TABLE attr_element_type (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_element_type'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_element_type__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: attr_element_val
--
-- Description: A table to hold attribute values.

CREATE TABLE attr_element_type_val (
    id           INTEGER      NOT NULL
                              DEFAULT NEXTVAL('seq_attr_element_type_val'),
    object__id   INTEGER      NOT NULL,
    attr__id     INTEGER      NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN      DEFAULT FALSE,
    active       BOOLEAN      NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_element_type_val__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_element_meta
--
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_element_type_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_element_type_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_element_type_meta__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--
CREATE UNIQUE INDEX udx_element_type__key_name ON element_type(LOWER(key_name));
CREATE INDEX fkx_et_type__element_type ON element_type(type__id);

CREATE INDEX fkx_class__element_type ON element_type(biz_class__id);

CREATE INDEX fkx_element_type__subelement__parent_id ON subelement_type(parent_id);
CREATE INDEX fkx_element_type__subelement__child_id ON subelement_type(child_id);
CREATE UNIQUE INDEX udx_subelement_type__parent__child ON subelement_type(parent_id, child_id);

CREATE UNIQUE INDEX udx_et_oc_id__et__oc_id ON element_type__output_channel(element_type__id, output_channel__id);
CREATE INDEX fkx_output_channel__et_oc ON element_type__output_channel(output_channel__id);
CREATE INDEX fkx_element__et_oc ON element_type__output_channel(element_type__id);

CREATE INDEX fkx_element_type__et_member ON element_type_member(object_id);
CREATE INDEX fkx_member__et_member ON element_type_member(member__id);

CREATE UNIQUE INDEX udx_element_type__site on element_type__site(element_type__id, site__id);

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_et__subsys__name ON attr_element_type(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_et__name ON attr_element_type(LOWER(name));
CREATE INDEX idx_attr_et__subsys ON attr_element_type(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_et_val__obj_attr ON attr_element_type_val (object__id, attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_et__attr_et_val ON attr_element_type_val(object__id);
CREATE INDEX fkx_attr_et__attr_et_val ON attr_element_type_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_et_meta__attr_name ON attr_element_type_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_et_meta__name ON attr_element_type_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_et__attr_et_meta ON attr_element_type_meta(attr__id);

-- FK index on element__site.
CREATE INDEX fkx_et__et__site__element_type__id ON element_type__site(element_type__id);
CREATE INDEX fkx_site__et__site__site__id ON element_type__site(site__id);
CREATE INDEX fkx_output_channel__et__site ON element_type__site(primary_oc__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- The sql to create the field_type table.
-- This maps to the Bric::ElementType::Parts::FieldType class.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the field_type table
CREATE SEQUENCE seq_field_type START 1024;

-- Unique IDs for the attr_field_type table
CREATE SEQUENCE seq_attr_field_type START 1024;

-- Unique IDs for each attr_field_type_val table
CREATE SEQUENCE seq_attr_field_type_val START 1024;

-- Unique IDs for the field_type_meta table
CREATE SEQUENCE seq_attr_field_type_meta START 1024;

-- -----------------------------------------------------------------------------
-- Table: field_type
--
-- Description: This is the table that contains the name and rules for fields
--         types. It contains references to the element_type table. The place
--         column represents the order that this is to be represented with in
--         its container.
--


CREATE TABLE field_type (
    id               INTEGER         NOT NULL
                                     DEFAULT NEXTVAL('seq_field_type'),
    element_type__id INTEGER         NOT NULL,
    name             TEXT            NOT NULL,
    key_name         TEXT            NOT NULL,
    description      TEXT,
    place            INTEGER         NOT NULL,
    min_occurrence   INTEGER         NOT NULL DEFAULT 0,
    max_occurrence   INTEGER         NOT NULL DEFAULT 0,
    autopopulated    BOOLEAN         NOT NULL DEFAULT FALSE,
    max_length       INTEGER         NOT NULL DEFAULT 0,
    sql_type         VARCHAR(30)     NOT NULL DEFAULT 'short',
    widget_type      VARCHAR(30)     NOT NULL DEFAULT 'text',
    precision        SMALLINT,
    cols             INTEGER         NOT NULL,
    rows             INTEGER         NOT NULL,
    length           INTEGER         NOT NULL,
    vals             TEXT,
    multiple         BOOLEAN         NOT NULL DEFAULT FALSE,
    default_val      TEXT,
    active           BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_field_type__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_field_type
--
-- Description: A table to represent types of attributes. A type is defined by
--              its subsystem, its field_type ID and an attribute name.

CREATE TABLE attr_field_type (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_field_type'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_field_type__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_field_type_val
--
-- Description: A table to hold attribute values.

CREATE TABLE attr_field_type_val (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_attr_field_type_val'),
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_field_type_val__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_field_type_meta
--
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_field_type_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_field_type_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     TEXT,
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_field_type_meta__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--

CREATE UNIQUE INDEX udx_field_type__key_name__et_id ON field_type(lower_text_num(key_name, element_type__id));
CREATE INDEX idx_field_type__name__at_id ON field_type(LOWER(name));
CREATE INDEX fkx_element_type__field_type on field_type(element_type__id);

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_field_type__subsys__name ON attr_field_type(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_field_type__name ON attr_field_type(LOWER(name));
CREATE INDEX idx_attr_field_type__subsys ON attr_field_type(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_field_type_val__obj_attr ON attr_field_type_val (object__id, attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_field_type__attr_field_type_val ON attr_field_type_val(object__id);
CREATE INDEX fkx_attr_field_type__attr_field_type_val ON attr_field_type_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_field_type_meta__attr_name ON attr_field_type_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_field_type_meta__name ON attr_field_type_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_field_type__attr_field_type_meta ON attr_field_type_meta(attr__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- This SQL creates the tables necessary for the keyword object.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the main keyword table
CREATE SEQUENCE seq_keyword START 1024;
CREATE SEQUENCE seq_keyword_member START 1024;


-- -----------------------------------------------------------------------------
-- Table: KEYWORD
--
-- Description: The main keyword table.

CREATE TABLE keyword (
    id               INTEGER       NOT NULL
                                   DEFAULT NEXTVAL('seq_keyword'),
    name             VARCHAR(256)  NOT NULL,
    screen_name      VARCHAR(256)  NOT NULL,
    sort_name        VARCHAR(256)  NOT NULL,
    active           BOOLEAN       NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_keyword__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: story_keyword
-- 
-- Description: The link between stories and keywords
--

CREATE TABLE story_keyword (
    story_id          INTEGER  NOT NULL,
    keyword_id        INTEGER  NOT NULL,
    CONSTRAINT pk_story_keyword PRIMARY KEY (story_id, keyword_id)
);


-- -----------------------------------------------------------------------------
-- Table: media_keyword
-- 
-- Description: The link between media and keywords
--

CREATE TABLE media_keyword (
    media_id         INTEGER  NOT NULL,
    keyword_id       INTEGER  NOT NULL,
    CONSTRAINT pk_media_keyword PRIMARY KEY (media_id, keyword_id)
);

-- -----------------------------------------------------------------------------
-- Table: category_keyword
-- 
-- Description: The link between categories and keywords
--

CREATE TABLE category_keyword (
    category_id       INTEGER  NOT NULL,
    keyword_id        INTEGER  NOT NULL,
    CONSTRAINT pk_category_keyword PRIMARY KEY (category_id, keyword_id)
);

--
-- TABLE: keyword_member
--

CREATE TABLE keyword_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_keyword_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_keyword_member__id PRIMARY KEY (id)
);




-- -----------------------------------------------------------------------------
-- Indexes

CREATE UNIQUE INDEX udx_keyword__name ON keyword(LOWER(name));
CREATE INDEX idx_keyword__screen_name ON keyword(LOWER(screen_name));
CREATE INDEX idx_keyword__sort_name   ON keyword(LOWER(sort_name));

CREATE INDEX fkx_keyword__keyword_member ON keyword_member(object_id);
CREATE INDEX fkx_member__keyword_member ON keyword_member(member__id);

CREATE INDEX fkx_keyword__story_keyword ON story_keyword(keyword_id);
CREATE INDEX fkx_story__story_keyword ON story_keyword(story_id);

CREATE INDEX fkx_keyword__media_keyword ON media_keyword(keyword_id);
CREATE INDEX fkx_media__media_keyword ON media_keyword(media_id);

CREATE INDEX fkx_keyword__category_keyword ON category_keyword(keyword_id);
CREATE INDEX fkx_category__category_keyword ON category_keyword(category_id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>


-- This DDL creates the table structure for Bric::Org objects.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_org START 1024;

-- 
-- TABLE: org 
--
CREATE TABLE org (
    id           INTEGER           NOT NULL
                                   DEFAULT NEXTVAL('seq_org'),
    name         VARCHAR(64)       NOT NULL,
    long_name    VARCHAR(128),
    personal     BOOLEAN           NOT NULL DEFAULT FALSE,
    active       BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_org__id PRIMARY KEY (id)
);



-- 
-- INDEXES.
--
CREATE INDEX idx_org__name ON org (LOWER(name));
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>


-- This DDL creates the table structure for Bric::Org::Parts::Address objects.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_addr START 1024;
CREATE SEQUENCE seq_addr_part START 1024;
CREATE SEQUENCE seq_addr_part_type START 1024;


-- 
-- TABLE: addr 
--

CREATE TABLE addr (
    id         INTEGER             NOT NULL
                                   DEFAULT NEXTVAL('seq_addr'),
    org__id    INTEGER             NOT NULL,
    type       VARCHAR(64),
    active     BOOLEAN             NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_addr__id PRIMARY KEY (id)
);


-- 
-- TABLE: addr_part_type 
--

CREATE TABLE addr_part_type (
    id         INTEGER             NOT NULL
                                   DEFAULT NEXTVAL('seq_addr_part_type'),
    name      VARCHAR(64)          NOT NULL,
    active     BOOLEAN             NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_addr_part_type__id PRIMARY KEY (id)
);


-- 
-- TABLE: addr_part
--

CREATE TABLE addr_part (
    id                    INTEGER         NOT NULL
                                          DEFAULT NEXTVAL('seq_addr_part'),
    addr__id              INTEGER         NOT NULL,
    addr_part_type__id    INTEGER         NOT NULL,
    value                 VARCHAR(256)    NOT NULL,
    CONSTRAINT pk_addr_part__id PRIMARY KEY (id)
);

-- 
-- TABLE: person_org__addr 
--

CREATE TABLE person_org__addr(
    addr__id          INTEGER           NOT NULL,
    person_org__id    INTEGER           NOT NULL,
    CONSTRAINT pk_person_org__addr__all PRIMARY KEY (addr__id,person_org__id)
);


--
-- INDEXES.
--
CREATE INDEX idx_addr__type ON addr(LOWER(type));
CREATE UNIQUE INDEX udx_addr_part_type__name ON addr_part_type(LOWER(name));
CREATE INDEX idx_addr_part__value ON addr_part(LOWER(value));

CREATE INDEX fkx_org__addr ON addr(org__id);
CREATE INDEX fkx_addr__addr_part ON addr_part(addr__id);
CREATE INDEX fkx_addr_part_type__addr_part ON addr_part(addr_part_type__id);
CREATE INDEX fkx_addr__person_org_addr ON person_org__addr(addr__id);
CREATE INDEX fk_person_org__pers_org_addr ON person_org__addr(person_org__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>

-- This DDL creates the table structure for Bric::Org::Person objects.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_person_org START 1024;

-- 
-- TABLE: person_org 
--
CREATE TABLE person_org(
    id            INTEGER           NOT NULL
                                    DEFAULT NEXTVAL('seq_person_org'),
    person__id    INTEGER           NOT NULL,
    org__id       INTEGER           NOT NULL,
    role          VARCHAR(64),
    department    VARCHAR(64),
    title         VARCHAR(64),
    active        BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_person_org__id PRIMARY KEY (id)
);


-- 
-- INDEXES.
--
CREATE UNIQUE INDEX udx_person_org__person__org ON person_org(person__id, org__id, role);
CREATE INDEX idx_person_org__department ON person_org(LOWER(department));
CREATE INDEX fkx_person__person_org ON person_org(person__id);
CREATE INDEX fkx_org__person_org ON person_org(org__id);
--
-- Project: Bricolage Business API
--
-- Author: David Wheeler <david@justatheory.com>


-- This DDL creates the table structure for Bric::BC::Org::Source objects.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_source START 1024;


-- 
-- TABLE: source
--

CREATE TABLE source (
    id            INTEGER           NOT NULL
                                    DEFAULT NEXTVAL('seq_source'),
    org__id       INTEGER           NOT NULL,
    name          VARCHAR(64)       NOT NULL,
    description   VARCHAR(256),
    expire        SMALLINT          NOT NULL DEFAULT 0,
    active        BOOLEAN           DEFAULT TRUE,
    CONSTRAINT pk_source__id PRIMARY KEY (id)
);


-- 
-- INDEXES.
--
CREATE UNIQUE INDEX udx_source_name ON source(name);
CREATE INDEX fkx_source__org on source(org__id);




-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- Description: The table that holds the registered Output Channels.
--                This maps to the Bric::OutputChannel Class.
--
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the output_channel table
CREATE SEQUENCE seq_output_channel START 1024;
CREATE SEQUENCE seq_output_channel_include START 1024;
CREATE SEQUENCE seq_output_channel_member START 1024;

-- -----------------------------------------------------------------------------
-- Table output_channel
--
-- Description: Holds info on the various output channels and is referenced
--                 by templates and elements
--
--

CREATE TABLE output_channel (
    id             INTEGER            NOT NULL
                                    DEFAULT NEXTVAL('seq_output_channel'),
    name             VARCHAR(64)    NOT NULL,
    description      VARCHAR(256),
    site__id         INTEGER        NOT NULL,
    protocol         VARCHAR(16),
    filename         VARCHAR(32)    NOT NULL,
    file_ext         VARCHAR(32),
    primary_ce       BOOLEAN,
    uri_format       VARCHAR(64)    NOT NULL,
    fixed_uri_format VARCHAR(64)    NOT NULL,
    uri_case         INT2           NOT NULL DEFAULT 1
                                     CONSTRAINT ck_output_channel__uri_case
                                     CHECK (uri_case IN (1,2,3)),
    use_slug         BOOLEAN        NOT NULL DEFAULT FALSE,
    burner           INT2           NOT NULL DEFAULT 1,
    active           BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_output_channel__id PRIMARY KEY (id)
);


--
-- TABLE: output_channel_include
--

CREATE TABLE output_channel_include (
    id                         INTEGER  NOT NULL
                               DEFAULT NEXTVAL('seq_output_channel_include'),
    output_channel__id         INTEGER  NOT NULL,
    include_oc_id              INTEGER  NOT NULL
                               CONSTRAINT ck_oc_include__include_oc_id
                                 CHECK (include_oc_id <> output_channel__id),
    CONSTRAINT pk_output_channel_include__id PRIMARY KEY (id)
);

--
-- TABLE: output_channel_member
--

CREATE TABLE output_channel_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_output_channel_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_output_channel_member__id PRIMARY KEY (id)
);

-- 
-- INDEXES.
--
CREATE UNIQUE INDEX udx_output_channel__name_site
ON output_channel(lower_text_num(name, site__id));

CREATE INDEX fkx_site__output_channel ON output_channel(site__id);
CREATE INDEX idx_output_channel__filename ON output_channel(LOWER(filename));
CREATE INDEX idx_output_channel__file_ext ON output_channel(LOWER(file_ext));

CREATE INDEX fkx_output_channel__oc_include ON output_channel_include(output_channel__id);
CREATE INDEX fkx_oc__oc_include_inc ON output_channel_include(include_oc_id);
CREATE UNIQUE INDEX udx_output_channel_include ON output_channel_include(output_channel__id, include_oc_id);
CREATE INDEX fkx_output_channel__oc_member ON output_channel_member(object_id);
CREATE INDEX fkx_member__oc_member ON output_channel_member(member__id);

-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>

-- This DDL creates the basic table for all Bric::Person objects. The indexes are
-- suggestions.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_person START 1024;
CREATE SEQUENCE seq_person_member START 1024;

-- 
-- TABLE: person 
--

CREATE TABLE person (
    id        INTEGER           NOT NULL
                                DEFAULT NEXTVAL('seq_person'),
    prefix    VARCHAR(32),
    lname     VARCHAR(64),
    fname     VARCHAR(64),
    mname     VARCHAR(64),
    suffix    VARCHAR(32),
    active    BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_person__id PRIMARY KEY (id)
);



--
-- TABLE: person_member
--

CREATE TABLE person_member (
    id          INTEGER         NOT NULL
                                DEFAULT NEXTVAL('seq_person_member'),
    object_id   INTEGER         NOT NULL,
    member__id  INTEGER         NOT NULL,
    CONSTRAINT pk_person_member__id PRIMARY KEY (id)
);


-- 
-- INDEXES.
--

CREATE INDEX idx_person__lname ON person(LOWER(lname));
CREATE INDEX idx_person__fname ON person(LOWER(fname));
CREATE INDEX idx_person__mname ON person(LOWER(mname));

CREATE INDEX fkx_person__person_member ON person_member(object_id);
CREATE INDEX fkx_member__person_member ON person_member(member__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.2
-- Author: David Wheeler <david@justatheory.com>

-- This DDL creates the basic table for Bric::Person::Usr objects, and
-- establishes its relationship with Bric::Person. The login field must be unique,
-- hence the udx_usr__login index.


-- 
-- TABLE: usr 
--

CREATE TABLE usr (
    id           INTEGER           NOT NULL,
    login        VARCHAR(128)      NOT NULL,
    password     CHAR(32)          NOT NULL,
    active       BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_usr__id PRIMARY KEY (id)
);

-- 
-- INDEXES.
--
CREATE INDEX idx_usr__login ON usr(LOWER(login));
CREATE UNIQUE INDEX udx_usr__login ON usr(LOWER(login)) WHERE active = '1';

-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>

-- This DDL creates the basic table for all Bric::Site objects. The indexes are
-- suggestions.

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_site_member START 1024;


-- 
-- TABLE: site 
--

CREATE TABLE site (
    id          INTEGER         NOT NULL,
    name        TEXT,
    description TEXT,
    domain_name TEXT,
    active      BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_site__id PRIMARY KEY (id)
);



--
-- TABLE: site_member
--

CREATE TABLE site_member (
    id          INTEGER  NOT NULL
                         DEFAULT NEXTVAL('seq_site_member'),
    object_id   INTEGER  NOT NULL,
    member__id  INTEGER  NOT NULL,
    CONSTRAINT pk_site_member__id PRIMARY KEY (id)
);


-- 
-- INDEXES.
--

CREATE UNIQUE INDEX udx_site__name ON site(LOWER(name));
CREATE UNIQUE INDEX udx_site__domain_name ON site(LOWER(domain_name));
CREATE INDEX fkx_site__site_member ON site_member(object_id);
CREATE INDEX fkx_member__site_member ON site_member(member__id);
-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- The database tables for the Workflow class.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- A sequence of unique IDs for the workflow table.
CREATE SEQUENCE seq_workflow START  1024;
CREATE SEQUENCE seq_workflow_member START 1024;

-- -----------------------------------------------------------------------------
-- Table: workflow
--
-- Description: The main workflow table.

CREATE TABLE workflow (
    id               INTEGER      NOT NULL
                                  DEFAULT NEXTVAL('seq_workflow'),
    name             VARCHAR(64)  NOT NULL,
    description      VARCHAR(256) NOT NULL,

    all_desk_grp_id  INTEGER      NOT NULL,
    req_desk_grp_id  INTEGER      NOT NULL,
    asset_grp_id     INTEGER      NOT NULL,
    head_desk_id     INTEGER      NOT NULL,
    type             INT2         NOT NULL
                                  DEFAULT 1
                                  CONSTRAINT ck_workflow__type
                                    CHECK (type IN (1,2,3)),
    active           BOOLEAN        NOT NULL DEFAULT TRUE,
    site__id         INTEGER      NOT NULL,
    CONSTRAINT pk_workflow__id PRIMARY KEY (id)
);

--
-- TABLE: workflow_member
--

CREATE TABLE workflow_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_workflow_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_workflow_member__id PRIMARY KEY (id)
);


-- 
-- INDEXES.
--

CREATE UNIQUE INDEX udx_workflow__name__site__id
ON workflow(lower_text_num(name, site__id));
CREATE INDEX fkx_site__workflow__site__id ON workflow(site__id);
CREATE INDEX fkx_grp__workflow__all_desk_grp_id ON workflow(all_desk_grp_id);
CREATE INDEX fkx_grp__workflow__req_desk_grp_id ON workflow(req_desk_grp_id);
CREATE INDEX fkx_grp__workflow__asset_grp_id ON workflow(asset_grp_id);
CREATE INDEX fkx_workflow__workflow_member ON workflow_member(object_id);
CREATE INDEX fkx_member__workflow_member ON workflow_member(member__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- The database tables for the Bric::Workflow::Parts::Desk class.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- A sequence of unique IDs for the desk table.
CREATE SEQUENCE seq_desk START  1024;

-- Unique IDs for the desk member ordering.
CREATE SEQUENCE seq_desk_member START  1024;

-- -----------------------------------------------------------------------------
-- Table: desk
--
-- Description: Represents a desk in the workflow

CREATE TABLE desk (
    id              INTEGER       NOT NULL
                                  DEFAULT NEXTVAL('seq_desk'),
    name            VARCHAR(64)   NOT NULL,
    description     VARCHAR(256),
    pre_chk_rules   INTEGER,
    post_chk_rules  INTEGER,
    asset_grp       INTEGER,
    publish         BOOLEAN       NOT NULL DEFAULT FALSE,
    active          BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_desk__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: desk_member
-- 
-- Description: The link between desk objects and member objects
--

CREATE TABLE desk_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_desk_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_desk_member__id PRIMARY KEY (id)
);


--
-- INDEXES.
--
CREATE UNIQUE INDEX udx_desk__name ON desk(LOWER(name));
CREATE INDEX fkx_asset_grp__desk ON desk(asset_grp);
CREATE INDEX fkx_pre_grp__desk ON desk(pre_chk_rules);
CREATE INDEX fkx_post_grp__desk ON desk(post_chk_rules);

CREATE INDEX fkx_desk__desk_member ON desk_member(object_id);
CREATE INDEX fkx_member__desk_member ON desk_member(member__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

--
-- SEQUENCES.
--
CREATE SEQUENCE seq_action START 1024;
CREATE SEQUENCE seq_attr_action START 1024;
CREATE SEQUENCE seq_attr_action_val START 1024;
CREATE SEQUENCE seq_attr_action_meta START 1024;


-- 
-- TABLE: action 
--

CREATE TABLE action (
    id               INTEGER           NOT NULL
                                       DEFAULT NEXTVAL('seq_action'),
    ord              INT2              NOT NULL,
    server_type__id  INTEGER           NOT NULL,
    action_type__id  INTEGER           NOT NULL,
    active           BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_action__id PRIMARY KEY (id)
);


-- 
-- TABLE: attr_action
--

CREATE TABLE attr_action (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_action'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_action__id PRIMARY KEY (id)
);


-- 
-- TABLE: attr_action_val
--

CREATE TABLE attr_action_val (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_attr_action_val'),
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_action_val__id PRIMARY KEY (id)
);


-- 
-- TABLE: attr_action_meta
--

CREATE TABLE attr_action_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_action_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_action_meta__id PRIMARY KEY (id)
);

-- 
-- INDEXES. 
--
CREATE INDEX fkx_action_type__action ON action(action_type__id);
CREATE INDEX fkx_server_type__action ON action(server_type__id);

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_action__subsys__name ON attr_action(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_action__name ON attr_action(LOWER(name));
CREATE INDEX idx_attr_action__subsys ON attr_action(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_action_val__obj_attr ON attr_action_val (object__id,attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_action__attr_action_val ON attr_action_val(object__id);
CREATE INDEX fkx_attr_action__attr_action_val ON attr_action_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_action_meta__attr_name ON attr_action_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_action_meta__name ON attr_action_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_action__attr_action_meta ON attr_action_meta(attr__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

--
-- SEQUENCES.
--
CREATE SEQUENCE seq_action_type START 1024;


-- 
-- TABLE: action_type 
--

CREATE TABLE action_type (
    id            INTEGER           NOT NULL
                                    DEFAULT NEXTVAL('seq_action_type'),
    name          VARCHAR(64)       NOT NULL,
    description   VARCHAR(256),
    active        BOOLEAN           NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_action_type__id PRIMARY KEY (id)
);


--
-- TABLE: action_type__media_type
--

CREATE TABLE action_type__media_type (
    action_type__id   INTEGER          NOT NULL,
    media_type__id    INTEGER          NOT NULL,
    CONSTRAINT pk_action__media_type PRIMARY KEY (action_type__id, media_type__id)
);


-- 
-- INDEXES. 
--

CREATE UNIQUE INDEX udx_action_type__name ON action_type(LOWER(name));
CREATE INDEX fkx_media_type__at_mt ON action_type__media_type(media_type__id);
CREATE INDEX fkx_action_type__at_mt ON action_type__media_type(action_type__id);



-- Project: Bricolage Business API
-- File:    Resource.sql
--
-- Author: David Wheeler <david@justatheory.com>
--


--
-- SEQUENCES.
--
CREATE SEQUENCE seq_resource START 1024;


-- 
-- TABLE: media__resource 
--

CREATE TABLE media__resource(
    resource__id    INTEGER           NOT NULL,
    media__id       INTEGER           NOT NULL,
    CONSTRAINT pk_media__resource PRIMARY KEY (media__id, resource__id)
);


-- 
-- TABLE: resource 
--

CREATE TABLE resource(
    id                  INTEGER           NOT NULL
                                          DEFAULT NEXTVAL('seq_resource'),
    parent_id           INTEGER,
    media_type__id      INTEGER           NOT NULL,
    path                VARCHAR(256)      NOT NULL,
    uri                 VARCHAR(256)      NOT NULL,
    size                INTEGER           NOT NULL,
    mod_time            TIMESTAMP         NOT NULL,
    is_dir              BOOLEAN           NOT NULL,
    CONSTRAINT pk_resource__id PRIMARY KEY (id)
);


-- 
-- TABLE: story__resource 
--

CREATE TABLE story__resource(
    story__id       INTEGER           NOT NULL,
    resource__id    INTEGER           NOT NULL,
    CONSTRAINT pk_story__resource PRIMARY KEY (story__id,resource__id)
);

-- 
-- INDEXES. 
--

CREATE UNIQUE INDEX udx_resource__path__uri ON resource(path, uri);
CREATE INDEX idx_resource__path ON resource(path);
CREATE INDEX idx_resource__uri ON resource(uri);
CREATE INDEX idx_resource__mod_time ON resource(mod_time);
CREATE INDEX fkx_media_type__resource ON resource(media_type__id);
CREATE INDEX fkx_resource__resource ON resource(parent_id);

CREATE INDEX fkx_resource__media__resource ON media__resource(resource__id);
CREATE INDEX fkx_media__media__resource ON media__resource(media__id);

CREATE INDEX fkx_story__story__resource ON story__resource(story__id);
CREATE INDEX fkx_resource__story__resource ON story__resource(resource__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- Sequences.
--
CREATE SEQUENCE seq_server START 1024;

-- 
-- TABLE: server 
--
CREATE TABLE server(
    id                 INTEGER          NOT NULL
                                        DEFAULT NEXTVAL('seq_server'),
    server_type__id    INTEGER          NOT NULL,
    host_name          VARCHAR(128)     NOT NULL,
    os                   CHAR(5)            NOT NULL,
    doc_root           VARCHAR(128)     NOT NULL,
    login              VARCHAR(64),
    password           VARCHAR(64),
    cookie             VARCHAR(512),
    active             BOOLEAN          NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_server__id PRIMARY KEY (id)
);


-- 
-- Indexes.
--
CREATE UNIQUE INDEX udx_server__name__st_id ON server(host_name, server_type__id);
CREATE INDEX fkx_server_type__server ON server(server_type__id);
CREATE INDEX idx_server__os ON server(os);


-- Project: Bricolage Business API
-- File:    ServerType.sql
--
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- Sequences.
--
CREATE SEQUENCE seq_server_type START 1024;
CREATE SEQUENCE seq_dest_member START 1024;

-- 
-- TABLE: server_type 
--

CREATE TABLE server_type(
    id             INTEGER           NOT NULL
                                     DEFAULT NEXTVAL('seq_server_type'),
    class__id      INTEGER           NOT NULL,
    name           VARCHAR(64)       NOT NULL,
    description    VARCHAR(256),
    site__id       INTEGER           NOT NULL,
    copyable       BOOLEAN           NOT NULL DEFAULT FALSE,
    publish        BOOLEAN           NOT NULL DEFAULT TRUE,
    preview        BOOLEAN           NOT NULL DEFAULT FALSE,
    active         BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_server_type__id PRIMARY KEY (id)
);


-- 
-- TABLE: server_type__output__channel
--

CREATE TABLE server_type__output_channel(
    server_type__id    INTEGER         NOT NULL,
    output_channel__id INTEGER         NOT NULL,
    CONSTRAINT pk_server_type__output_channel
      PRIMARY KEY (server_type__id, output_channel__id)
);


--
-- TABLE: dest_member
--

CREATE TABLE dest_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_dest_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_dest_member__id PRIMARY KEY (id)
);


-- 
-- Indexes.
--
CREATE UNIQUE INDEX udx_server_type__name_site
ON server_type(lower_text_num(name, site__id));

CREATE INDEX fkx_site__server_type ON server_type(site__id);
CREATE INDEX fkx_class__server_type ON server_type(class__id);

CREATE INDEX fkx_server_type__st_oc ON server_type__output_channel(server_type__id);
CREATE INDEX fk_output_channel__st_oc ON server_type__output_channel(output_channel__id);

CREATE INDEX fkx_dest__dest_member ON dest_member(object_id);
CREATE INDEX fkx_member__dest_member ON dest_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_alert START 1024;


-- 
-- TABLE: alert 
--

CREATE TABLE alert(
    id                INTEGER           NOT NULL
                                        DEFAULT NEXTVAL('seq_alert'),
    alert_type__id    INTEGER           NOT NULL,
    event__id         INTEGER           NOT NULL,
    subject           VARCHAR(128),
    message           VARCHAR(512),
    ttimestamp         TIMESTAMP         NOT NULL
                                        DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_alert__id PRIMARY KEY (id)
);

-- 
-- INDEXES.
--
 
CREATE INDEX idx_alert__timestamp ON alert(timestamp);
CREATE INDEX fkx_alert_type__alert ON alert(alert_type__id);
CREATE INDEX fkx_event__alert ON alert(event__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_alert_type START 1024;


-- 
-- TABLE: alert_type 
--

CREATE TABLE alert_type (
    id                INTEGER           NOT NULL
                                        DEFAULT NEXTVAL('seq_alert_type'),
    event_type__id    INTEGER           NOT NULL,
    usr__id           INTEGER           NOT NULL,
    name              VARCHAR(64)       NOT NULL,
    subject           VARCHAR(128),
    message           VARCHAR(512),
    active            BOOLEAN           NOT NULL DEFAULT TRUE,
    del               BOOLEAN           NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_alert_type__id PRIMARY KEY (id)
);


-- 
-- TABLE: alert_type__grp__contact 
--

CREATE TABLE alert_type__grp__contact(
    alert_type__id    INTEGER           NOT NULL,
    contact__id       INTEGER           NOT NULL,
    grp__id           INTEGER           NOT NULL,
    CONSTRAINT pk_alert_type__grp__contact PRIMARY KEY (alert_type__id, contact__id, grp__id)
);


-- 
-- TABLE: alert_type__usr__contact 
--

CREATE TABLE alert_type__usr__contact(
    alert_type__id    INTEGER           NOT NULL,
    contact__id       INTEGER           NOT NULL,
    usr__id           INTEGER           NOT NULL,
    CONSTRAINT pk_alert_type__usr__contact PRIMARY KEY (alert_type__id, usr__id, contact__id)
);


-- 
-- INDEXES.
--

-- alert_type
CREATE UNIQUE INDEX udx_alert_type__name__usr__id
ON alert_type(lower_text_num(name, usr__id));

CREATE INDEX idx_alert_type__name ON alert_type(LOWER(name));
CREATE INDEX fkx_event_type__alert_type ON alert_type(event_type__id);
CREATE INDEX fkx_usr__alert_type ON alert_type(usr__id);

-- alert_type__grp__contact
CREATE INDEX fkx_alert_type__grp__contact ON alert_type__grp__contact(alert_type__id);
CREATE INDEX fkx_contact__grp__contact ON alert_type__grp__contact(contact__id);
CREATE INDEX fkx_grp__grp__contact ON alert_type__grp__contact(grp__id);

-- alert_type__usr__contact
CREATE INDEX fkx_alert_type__at_user__cont ON alert_type__usr__contact(alert_type__id);
CREATE INDEX fkx_contact__at_usr__contact  ON alert_type__usr__contact(contact__id);
CREATE INDEX fkx_usr__at_usr__contact ON alert_type__usr__contact(usr__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_alert_type_rule START  1024;

-- 
-- TABLE: alert_type_rule
--

CREATE TABLE alert_type_rule(
    id                INTEGER         NOT NULL
                                      DEFAULT NEXTVAL('seq_alert_type_rule'),
    alert_type__id    INTEGER         NOT NULL,
    attr              VARCHAR(64)     NOT NULL,
    operator          CHAR(2)         NOT NULL,
    value             VARCHAR(256)    NOT NULL,
    CONSTRAINT pk_alert_type_rule__id PRIMARY KEY (id)
);

-- 
-- INDEXS.
--

CREATE INDEX idx_alert_type_rule__attr ON alert_type_rule(LOWER(attr));
CREATE INDEX idx_alert_type_rule__value ON alert_type_rule(LOWER(value));
CREATE INDEX fkx_alert_type__at_rule ON alert_type_rule(alert_type__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_alerted START  1024;

-- 
-- TABLE: alerted 
--

CREATE TABLE alerted(
    id           INTEGER           NOT NULL
                                   DEFAULT NEXTVAL('seq_alerted'),
    usr__id      INTEGER           NOT NULL,
    alert__id    INTEGER           NOT NULL,
    ack_time     TIMESTAMP,
    CONSTRAINT pk_alerted__id PRIMARY KEY (id)
);


-- 
-- TABLE: alerted__contact_value 
--

CREATE TABLE alerted__contact_value(
    alerted__id                INTEGER         NOT NULL,
    contact__id             INTEGER         NOT NULL,
    contact_value__value    VARCHAR(256)    NOT NULL,
    sent_time               TIMESTAMP,
    CONSTRAINT pk_alerted__contact_value PRIMARY KEY (alerted__id, contact__id, contact_value__value)
);

-- 
-- INDEXES.
--

-- alerted
CREATE INDEX idx_alerted__ack_time ON alerted(ack_time);
CREATE INDEX fkx_alert__alerted ON alerted(alert__id);
CREATE INDEX fkx_usr__alerted ON alerted(usr__id);

-- alerted__contact_value
CREATE INDEX idx_ac_value__sent_time ON alerted__contact_value(sent_time);
CREATE INDEX idx_ac_value__cv__value ON alerted__contact_value(contact_value__value);
CREATE INDEX fkx_alerted__alerted__contact ON alerted__contact_value(alerted__id);
CREATE INDEX fkx_contact__alerted__cont ON alerted__contact_value(contact__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- -----------------------------------------------------------------------------
-- Attribute.sql
--
--
-- This SQL creates the tables necessary for the attribute object.  This file
-- applies to attributes on the Bric::Person class.  Any other classes that 
-- require attributes need only duplicate these tables, changing 'person' to 
-- the correct class name.  Class names may be shortened to ensure that the
-- resulting table names are under the oracle 30 character name limit as long
-- as the resulting shortened class name is unique.
--

/* Commented out because attr_person won't be used in production (in this version).
   However, the examples still apply. --David, 23 Feb 2001

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the attr_person table
CREATE SEQUENCE seq_attr_person START 1024;

-- Unique IDs for each attr_person_*_val table
CREATE SEQUENCE seq_attr_person_val START 1024;

-- Unique IDs for the person_meta table
CREATE SEQUENCE seq_attr_person_meta START 1024;

-- -----------------------------------------------------------------------------
-- Table: attr_person
--
-- Description: A table to represent types of attributes.  A type is defined by
--              its subsystem, its person ID and an attribute name.

CREATE TABLE attr_person (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_person'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_person__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Table: attr_person_val
--
-- Description: A table to hold attribute values.

CREATE TABLE attr_person_val (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_attr_person_val'),
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_person_val__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Table: attr_person_meta
--
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_person_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_person_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_person_meta__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Indexes.
--

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_person__subsys__name ON attr_person(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_person__name ON attr_person(LOWER(name));
CREATE INDEX idx_attr_person__subsys ON attr_person(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_person_val__obj_attr ON attr_person_val (object__id,attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_person__attr_person_val ON attr_person_val(object__id);
CREATE INDEX fkx_attr_person__attr_person_val ON attr_person_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_person_meta__attr_name ON attr_person_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_person_meta__name ON attr_person_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_person__attr_person_meta ON attr_person_meta(attr__id);

*/


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Garth Webb <garth@perijove.com>
--
-- This is the SQL that will create the class table.
--

-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the class table
CREATE SEQUENCE seq_class START 1024;

-- -----------------------------------------------------------------------------
-- TABLE: class 
--        For keeping track of Perl classes.

CREATE TABLE class(
    id              INTEGER         NOT NULL
                                    DEFAULT NEXTVAL('seq_class'),
    key_name        VARCHAR(32)     NOT NULL
                                    CONSTRAINT ck_class__key_name
                                    CHECK (LOWER(key_name) = key_name),
    pkg_name        VARCHAR(128)    NOT NULL,
    disp_name       VARCHAR(128)    NOT NULL,
    plural_name        VARCHAR(128)    NOT NULL,
    description     VARCHAR(256),
    distributor     BOOLEAN         NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_class__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--

CREATE UNIQUE INDEX udx_class__key_name ON class(LOWER(key_name));
CREATE UNIQUE INDEX udx_class__pkg_name ON class(LOWER(pkg_name));
CREATE UNIQUE INDEX udx_class__disp__name ON class(LOWER(disp_name));

-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>


-- This DDL creates the basic tables for Bric::Util::Event objects.

-- 
-- SEQUENCES.
--
CREATE SEQUENCE seq_event START 1024;
CREATE SEQUENCE seq_event_attr START 1024;

-- 
-- TABLE: event 
--

CREATE TABLE event (
    id                INTEGER           NOT NULL
                                        DEFAULT NEXTVAL('seq_event'),
    event_type__id    INTEGER           NOT NULL,
    usr__id           INTEGER           NOT NULL,
    obj_id            INTEGER           NOT NULL,
    ttimestamp         TIMESTAMP         NOT NULL
                                        DEFAULT 'timeofday()::timestamptz',
    CONSTRAINT pk_event__id PRIMARY KEY (id)
);

-- 
-- TABLE: event_attr
--
CREATE TABLE event_attr (
    id                   INTEGER          NOT NULL
                                          DEFAULT NEXTVAL('seq_event_attr'),
    event__id            INTEGER          NOT NULL,
    event_type_attr__id  INTEGER          NOT NULL,
    value                VARCHAR(128),
    CONSTRAINT pk_event_attr__id PRIMARY KEY (id)
);

--
-- INDEXES.
--

CREATE INDEX fkx_event_type__event ON event(event_type__id);
CREATE INDEX fkx_usr__event ON event(usr__id);
CREATE INDEX idx_event__timestamp ON event(timestamp);
CREATE INDEX idx_event__obj_id ON event(obj_id);

CREATE INDEX fkx_event__event_attr ON event_attr(event__id);
CREATE INDEX fkx_event_type_attr__event_attr ON event_attr(event_type_attr__id);


--
-- ER/Studio 4.0 SQL Code Generation
-- Project:      Bricolage Business API
--
-- Target DBMS : Oracle 8
-- Author: David Wheeler <david@justatheory.com>

-- This DDL creates the basic table for all Bric::Util::EventType objects. It's
-- pretty easy - they're really just all groups.

--
-- SEQUENCES.
--
CREATE SEQUENCE seq_event_type START 1024;
CREATE SEQUENCE seq_event_type_attr START 1024;


-- 
-- TABLE: event_type
--

CREATE TABLE event_type (
    id              INTEGER         NOT NULL
                                    DEFAULT NEXTVAL('seq_event_type'),
    key_name        VARCHAR(64)     NOT NULL,
    name            VARCHAR(64)     NOT NULL,
    description     VARCHAR(256)    NOT NULL,
    class__id       INTEGER         NOT NULL,
    active          BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_event_type__id PRIMARY KEY (id)
);

-- 
-- TABLE: event_type_attr
--

CREATE TABLE event_type_attr (
    id              INTEGER         NOT NULL
                                    DEFAULT NEXTVAL('seq_event_type_attr'),
    event_type__id  INTEGER         NOT NULL,
    name            VARCHAR(64)     NOT NULL,
    CONSTRAINT pk_event_type_attr__id PRIMARY KEY (id)
);    


-- 
-- INDEXES.
--

CREATE UNIQUE INDEX udx_event_type__key_name ON event_type(LOWER(key_name));
CREATE UNIQUE INDEX udx_event_type__class_id__name ON event_type(class__id, name);

CREATE INDEX fkx_event_type__event_type_attr ON event_type_attr(event_type__id);

CREATE INDEX fkx_class__event_type ON event_type(class__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--

-- ----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the grp table
CREATE SEQUENCE seq_grp START  1024; 

-- ----------------------------------------------------------------------------
-- Table grp
-- 
-- Description: The grp table   Contains the name and description of the
--                 group and its parent if it has one
--

CREATE TABLE grp (
    id           INTEGER          NOT NULL
                                  DEFAULT NEXTVAL('seq_grp'),
    parent_id    INTEGER          CONSTRAINT ck_grp__parent_id_not_eq_id
                                    CHECK (parent_id <> id),
    class__id    INTEGER          NOT NULL,
    name         VARCHAR(64),
    description  VARCHAR(256),
    secret       BOOLEAN          NOT NULL DEFAULT TRUE,
    permanent    BOOLEAN          NOT NULL DEFAULT FALSE,
    active       BOOLEAN          NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_grp__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX idx_grp__name ON grp(LOWER(name));
CREATE INDEX idx_grp__description ON grp(LOWER(description));
CREATE INDEX fkx_grp__grp ON grp(parent_id);
CREATE INDEX fkx_class__grp ON grp(class__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_alert_type_member START 1024;

-- 
-- TABLE: alert_type_member 
--

CREATE TABLE alert_type_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_alert_type_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_alert_type_member__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_alert_type__alert_type_member ON alert_type_member(object_id);
CREATE INDEX fkx_member__alert_type_member ON alert_type_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>



CREATE SEQUENCE seq_contrib_type_member START 1024;


--
-- TABLE: contrib_type_member
--

CREATE TABLE contrib_type_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_contrib_type_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_contrib_type_member__id PRIMARY KEY (id)
);






CREATE INDEX fkx_contrib_type__ctype_member ON contrib_type_member(object_id);
CREATE INDEX fkx_member__ctype_member ON contrib_type_member(member__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>

-- 
-- TABLE: event_member
--


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_org_member START 1024;

-- 
-- TABLE: org_member 
--

CREATE TABLE org_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_org_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_org_member__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_org__org_member ON org_member(object_id);
CREATE INDEX fkx_member__org_member ON org_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- -----------------------------------------------------------------------------
-- Member.sql
-- 
--
-- The member table and the tables that map member back to their respective 
-- objects. The member table contains an id and a group id. The table that 
-- maps the object to its member contains an id an object id and a member id
--
-- Thought should be given to:
--         Ensuring that an object is not placed with in the same group twice
--        Making sure that an object that is deactivated from a group that is 
--            then put back in again will behave properly
--

-- -----------------------------------------------------------------------------
-- Sequences


-- Unique IDs for the grp_member table
CREATE SEQUENCE seq_grp_member START  1024;


-- -----------------------------------------------------------------------------
-- Table: grp_member
-- 
-- Description: The link between stroy objects and member objects
--

CREATE TABLE grp_member (
    id            INTEGER         NOT NULL
                                  DEFAULT NEXTVAL('seq_grp_member'),
    object_id     INTEGER         NOT NULL,
    member__id    INTEGER            NOT NULL,
    CONSTRAINT pk_grp_member__id PRIMARY KEY (id)
);


--
-- INDEXES
--

CREATE INDEX fkx_grp__grp_member ON grp_member(object_id);
CREATE INDEX fkx_member__grp_member ON grp_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- -----------------------------------------------------------------------------
-- Member.sql
-- 
--
-- The member table and the tables that map member bac to their respective 
-- objects.   The member table contains an id and a group id.   The table that 
-- maps the object to its member contains an id an object id and a member id
--
-- Thought should be given to:
--         Ensuring that an object is not placed with in the same group twice
--        Making sure that an object that is deactivated from a group that is 
--            then put back in again will behave properly
--

-- -----------------------------------------------------------------------------

-- 
-- SEQUENCES.
--
CREATE SEQUENCE seq_member START  1024;
CREATE SEQUENCE seq_story_member START 1024;

-- -----------------------------------------------------------------------------
-- Table: member
--
-- Description:    The table that creates a member of a group.   The obj_member 
-- table then links the objects to the member table
--

CREATE TABLE member (
    id         INTEGER        NOT NULL
                              DEFAULT NEXTVAL('seq_member'),
    grp__id    INTEGER        NOT NULL,
    class__id  INTEGER        NOT NULL,
    active     BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_member__id PRIMARY KEY (id)
);


--
-- INDEXES.
--
CREATE INDEX fkx_grp__member ON member(grp__id);
CREATE INDEX fkx_grp__class ON member(class__id);

-- Use the below section as an example to create new member tables for
-- other objects.
-- -----------------------------------------------------------------------------
-- Table: story_member
-- 
-- Description: The link between story objects and member objects
--

CREATE TABLE story_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_story_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_story_member__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_story__story_member ON story_member(object_id);
CREATE INDEX fkx_member__story_member ON story_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- -----------------------------------------------------------------------------
--
-- This SQL creates the tables necessary for the attribute object.  This file
-- applies to attributes on the Bric::Util::Grp class.  Any other classes that 
-- require attributes need only duplicate these tables, changing 'member' to 
-- the correct class name.  Class names may be shortened to ensure that the
-- resulting table names are under the oracle 30 character name limit as long
-- as the resulting shortened class name is unique.
--

-- -------------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the attr_member table
CREATE SEQUENCE seq_attr_member START  1024;

-- Unique IDs for each attr_member_*_val table
CREATE SEQUENCE seq_attr_member_val START  1024;

-- Unique IDs for the member_meta table
CREATE SEQUENCE seq_attr_member_meta START  1024;

-- Table: attr_member
--
-- Description: A table to represent types of attributes.  A type is defined by
--              its subsystem, its member ID and an attribute name.

CREATE TABLE attr_member (
    id         INTEGER       NOT NULL,
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_member__id PRIMARY KEY (id)
);

-- -------------------------------------------------------------------------------
-- Table: attr_member_val
-- Description: A table to hold attribute values.

CREATE TABLE attr_member_val (
    id           INTEGER         NOT NULL,
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_member_val__id PRIMARY KEY (id)
);

-- -------------------------------------------------------------------------------
-- Table: attr_member_meta
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_member_meta (
    id        INTEGER         NOT NULL,
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_member_meta__id PRIMARY KEY (id)
);

-- -----------------------------------------------------------------------------
-- Indexes.
--

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_member__subsys__name ON attr_member(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_member__name ON attr_member(LOWER(name));
CREATE INDEX idx_attr_member__subsys ON attr_member(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_member_val__obj_attr ON attr_member_val (object__id,attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_member__attr_member_val ON attr_member_val(object__id);
CREATE INDEX fkx_attr_member__attr_member_val ON attr_member_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_member_meta__attr_name ON attr_member_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_member_meta__name ON attr_member_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_member__attr_member_meta ON attr_member_meta(attr__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_source_member START 1024;

-- 
-- TABLE: source_member 
--

CREATE TABLE source_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_source_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_source_member__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_source__source_member ON source_member(object_id);
CREATE INDEX fkx_member__source_member ON source_member(member__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

-- 
-- SEQUENCES.
--

CREATE SEQUENCE seq_user_member START 1024;

-- 
-- TABLE: user_member 
--

CREATE TABLE user_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_user_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_user_member__id PRIMARY KEY (id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_user__user_member ON user_member(object_id);
CREATE INDEX fkx_member__user_member ON user_member(member__id);



-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--
-- -----------------------------------------------------------------------------
--
-- This SQL creates the tables necessary for the attribute object.  This file
-- applies to attributes on the Bric::Util::Grp class.  Any other classes that 
-- require attributes need only duplicate these tables, changing 'grp' to 
-- the correct class name.  Class names may be shortened to ensure that the
-- resulting table names are under the oracle 30 character name limit as long
-- as the resulting shortened class name is unique.
--

-- -------------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the attr_grp table
CREATE SEQUENCE seq_attr_grp START  1024;

-- Unique IDs for each attr_grp_*_val table
CREATE SEQUENCE seq_attr_grp_val START  1024;

-- Unique IDs for the grp_meta table
CREATE SEQUENCE seq_attr_grp_meta START  1024;

-- Table: attr_grp
--
-- Description: A table to represent types of attributes.  A type is defined by
--              its subsystem, its grp ID and an attribute name.

CREATE TABLE attr_grp (
    id         INTEGER       NOT NULL
                             DEFAULT NEXTVAL('seq_attr_grp'),
    subsys     VARCHAR(256)  NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    sql_type   VARCHAR(30)   NOT NULL,
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_grp__id PRIMARY KEY (id)
);

-- -------------------------------------------------------------------------------
-- Table: attr_grp_val
--
-- Description: A table to hold attribute values.


CREATE TABLE attr_grp_val (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_attr_grp_val'),
    object__id   INTEGER         NOT NULL,
    attr__id     INTEGER         NOT NULL,
    date_val     TIMESTAMP,
    short_val    VARCHAR(1024),
    blob_val     TEXT,
    serial       BOOLEAN         DEFAULT FALSE,
    active       BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_attr_grp_val__id PRIMARY KEY (id)
);


-- -------------------------------------------------------------------------------
-- Table: attr_grp_meta
--
-- Description: A table to represent metadata on types of attributes.

CREATE TABLE attr_grp_meta (
    id        INTEGER         NOT NULL
                              DEFAULT NEXTVAL('seq_attr_grp_meta'),
    attr__id  INTEGER         NOT NULL,
    name      VARCHAR(256)    NOT NULL,
    value     VARCHAR(2048),
    active    BOOLEAN         NOT NULL DEFAULT TRUE,
   CONSTRAINT pk_attr_grp_meta__id PRIMARY KEY (id)
);


-- -----------------------------------------------------------------------------
-- Indexes.
--

-- Unique index on subsystem/name pair
CREATE UNIQUE INDEX udx_attr_grp__subsys__name ON attr_grp(subsys, name);

-- Indexes on name and subsys.
CREATE INDEX idx_attr_grp__name ON attr_grp(LOWER(name));
CREATE INDEX idx_attr_grp__subsys ON attr_grp(LOWER(subsys));

-- Unique index on object__id/attr__id pair
CREATE UNIQUE INDEX udx_attr_grp_val__obj_attr ON attr_grp_val (object__id,attr__id);

-- FK indexes on object__id and attr__id.
CREATE INDEX fkx_grp__attr_grp_val ON attr_grp_val(object__id);
CREATE INDEX fkx_attr_grp__attr_grp_val ON attr_grp_val(attr__id);

-- Unique index on attr__id/name pair
CREATE UNIQUE INDEX udx_attr_grp_meta__attr_name ON attr_grp_meta (attr__id, name);

-- Index on meta name.
CREATE INDEX idx_attr_grp_meta__name ON attr_grp_meta(LOWER(name));

-- FK index on attr__id.
CREATE INDEX fkx_attr_grp__attr_grp_meta ON attr_grp_meta(attr__id);


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

--
-- SEQUENCES.
--
CREATE SEQUENCE seq_job START 1024;
CREATE SEQUENCE seq_job_member START 1024;

-- 
-- TABLE: job 
--

CREATE TABLE job (
    id                 INTEGER           NOT NULL
                                         DEFAULT NEXTVAL('seq_job'),
    name               TEXT              NOT NULL,
    usr__id            INTEGER           NOT NULL,
    sched_time         TIMESTAMP         NOT NULL
                                         DEFAULT CURRENT_TIMESTAMP,
    priority           INT2              NOT NULL 
                                         DEFAULT 3
                                         CONSTRAINT ck_job__priority 
                                           CHECK (priority BETWEEN 1 AND 5),
    comp_time          TIMESTAMP,
    expire             BOOLEAN           NOT NULL DEFAULT FALSE,
    failed             BOOLEAN           NOT NULL DEFAULT FALSE,
    tries              INT2              NOT NULL DEFAULT 0
                                         CONSTRAINT ck_job__tries
                                           CHECK (tries BETWEEN 0 AND 10),
    executing          BOOLEAN           NOT NULL DEFAULT FALSE,
    class__id          INTEGER           NOT NULL,
    story_instance__id INTEGER,
    media_instance__id INTEGER,
    error_message TEXT,
    CONSTRAINT pk_job__id PRIMARY KEY (id)
);


-- 
-- TABLE: job__resource 
--

CREATE TABLE job__resource(
    job__id         INTEGER           NOT NULL,
    resource__id    INTEGER           NOT NULL,
    CONSTRAINT pk_job__resource PRIMARY KEY (job__id,resource__id)
);


-- 
-- TABLE: job__server_type 
--

CREATE TABLE job__server_type(
    job__id             INTEGER        NOT NULL,
    server_type__id     INTEGER        NOT NULL,
    CONSTRAINT pk_job__server_type PRIMARY KEY (job__id,server_type__id)
);

--
-- TABLE: job_member
--

CREATE TABLE job_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_job_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_job_member__id PRIMARY KEY (id)
);


-- 
-- INDEXES. 
--
CREATE INDEX idx_job__name ON job(LOWER(name));
CREATE INDEX idx_job__sched_time ON job(sched_time);
CREATE INDEX idx_job__comp_time__is_null ON job(comp_time) WHERE comp_time is NULL;
CREATE INDEX idx_job__comp_time ON job(comp_time);
CREATE INDEX idx_job__executing ON job(executing);

CREATE INDEX fkx_story_instance__job ON job(story_instance__id)
WHERE story_instance__id IS NOT NULL;
CREATE INDEX fkx_media_instance__job ON job(media_instance__id)
WHERE media_instance__id IS NOT NULL;

CREATE INDEX fkx_job__job__resource ON job__resource(job__id);
CREATE INDEX fkx_usr__job ON job (usr__id);
CREATE INDEX fkx_resource__job__resource ON job__resource(resource__id);
CREATE INDEX fkx_job__job__server_type ON job__server_type(job__id);
CREATE INDEX fkx_srvr_type__job__srvr_type ON job__server_type(server_type__id);

CREATE INDEX fkx_job__job_member ON job_member(object_id);
CREATE INDEX fkx_member__job_member ON job_member(member__id);

-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: Michael Soderstrom <miraso@pacbell.net>
--

/* Commented out because we're not using language stuff at this point.
   By David.


-- -----------------------------------------------------------------------------
-- Sequences

-- Unique IDs for the language table
-- IDs under 1024 will contain dead languages
CREATE SEQUENCE seq_language START 1024;


-- -----------------------------------------------------------------------------
-- Table: language
--
-- Description: name and description of languages 
--              

CREATE TABLE language (
    id           INTEGER          NOT NULL
                                DEFAULT NEXTVAL('seq_language'),
    name         VARCHAR(64),
    description  VARCHAR(256),
    active       BOOLEAN        NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_language__id PRIMARY KEY (id)
);

CREATE UNIQUE INDEX udx_language__name ON language(LOWER(name));

*/


-- Project: Bricolage
--
-- Target DBMS: PostgreSQL 7.1.2
-- Author: David Wheeler <david@justatheory.com>
--

--
-- SEQUENCES.
--
CREATE SEQUENCE seq_media_type START 1024;
CREATE SEQUENCE seq_media_type_ext START 1024;
CREATE SEQUENCE seq_media_type_member START 1024;

-- 
-- TABLE: media_type 
--

CREATE TABLE media_type (
    id             INTEGER           NOT NULL
                                     DEFAULT NEXTVAL('seq_media_type'),
    name           VARCHAR(128)      NOT NULL,
    description    VARCHAR(256),
    active         BOOLEAN           NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_media_type__id PRIMARY KEY (id)
);


-- 
-- TABLE: media_type_ext
--

CREATE TABLE media_type_ext (
    id                  INTEGER           NOT NULL
                                          DEFAULT NEXTVAL('seq_media_type_ext'),
    media_type__id      INTEGER           NOT NULL,
    extension           VARCHAR(10)       NOT NULL,
    CONSTRAINT pk_media_type_ext__id PRIMARY KEY (id)
);


--
-- TABLE: media_type_member
--

CREATE TABLE media_type_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_media_type_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_media_type_member__id PRIMARY KEY (id)
);


-- 
-- INDEXES. 
--

CREATE UNIQUE INDEX udx_media_type__name ON media_type(LOWER(name));
CREATE UNIQUE INDEX udx_media_type_ext__extension ON media_type_ext(LOWER(extension));
CREATE INDEX fkx_media_type__media_type_ext ON media_type_ext(media_type__id);
CREATE INDEX fkx_media_type__media_type_member ON media_type_member(object_id);
CREATE INDEX fkx_member__media_type_member ON media_type_member(member__id);


--
-- Project: Bricolage API
--
-- Author: David Wheeler <david@justatheory.com>


--
-- SEQUENCES
--

CREATE SEQUENCE seq_pref START 1024;
CREATE SEQUENCE seq_usr_pref START 1024;
CREATE SEQUENCE seq_pref_member START 1024;



-- 
-- TABLE: pref
--        Global preferences.

CREATE TABLE pref (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_pref'),
    name         VARCHAR(64)     NOT NULL,
    description  VARCHAR(256),
    value        VARCHAR(256),
    def          VARCHAR(256),
    manual         BOOLEAN         NOT NULL DEFAULT FALSE,
    opt_type     VARCHAR(16)     NOT NULL,
    can_be_overridden  BOOLEAN   NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_pref__id PRIMARY KEY (id)
);

-- 
-- TABLE: usr_pref
--        Preferences overridden by a specific usr.

CREATE TABLE usr_pref (
    id           INTEGER         NOT NULL
                                 DEFAULT NEXTVAL('seq_usr_pref'),
    pref__id     INTEGER         NOT NULL,
    usr__id      INTEGER         NOT NULL,
    value        VARCHAR(256)    NOT NULL,
    CONSTRAINT pk_usr_pref__pref__id__value PRIMARY KEY (id)
);

-- 
-- TABLE: pref
--        Preference options.

CREATE TABLE pref_opt (
    pref__id     INTEGER         NOT NULL,
    value        VARCHAR(256)    NOT NULL,
    description  VARCHAR(256),
    CONSTRAINT pk_pref_opt__pref__id__value PRIMARY KEY (pref__id, value)
);


--
-- TABLE: pref_member
--

CREATE TABLE pref_member (
    id          INTEGER        NOT NULL
                               DEFAULT NEXTVAL('seq_pref_member'),
    object_id   INTEGER        NOT NULL,
    member__id  INTEGER        NOT NULL,
    CONSTRAINT pk_pref_member__id PRIMARY KEY (id)
);



-- 
-- INDEXES.
--

CREATE UNIQUE INDEX udx_pref__name ON pref(LOWER(name));
CREATE UNIQUE INDEX udx_usr_pref__pref__id__usr__id ON usr_pref(pref__id, usr__id);
CREATE INDEX idx_usr_pref__usr__id ON usr_pref(usr__id);
CREATE INDEX fkx_pref__pref__opt ON pref_opt(pref__id);
CREATE INDEX fkx_pref__pref_member ON pref_member(object_id);
CREATE INDEX fkx_member__pref_member ON pref_member(member__id);



--
-- Project: Bricolage API
--
-- Author: David Wheeler <david@justatheory.com>

--
-- SEQUENCES
--
-- Use in both grp_priv and usr_priv (if we ever implement the latter).
CREATE SEQUENCE seq_priv START 1024;


-- 
-- TABLE: grp_priv 
--        Privileges granted to user groups.

CREATE TABLE grp_priv (
    id         INTEGER           NOT NULL
                                 DEFAULT NEXTVAL('seq_priv'),
    grp__id    INTEGER           NOT NULL,
    value      INT2              NOT NULL
                                 CONSTRAINT ck_grp_priv__value
                                   CHECK (value BETWEEN 1 AND 255),
    mtime      TIMESTAMP         NOT NULL
                                 DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_grp_priv__id PRIMARY KEY (id)
);


-- 
-- TABLE: grp_priv__grp_member 
--        Ties group privileges to groups for whose members the privilege
--        is granted.

CREATE TABLE grp_priv__grp_member (
    grp_priv__id    INTEGER           NOT NULL,
    grp__id         INTEGER           NOT NULL,
    CONSTRAINT pk_grp_priv__grp_member PRIMARY KEY (grp_priv__id,grp__id)
);

--
-- INDEXES.
--
CREATE INDEX fkx_grp__grp_priv ON grp_priv(grp__id);
CREATE INDEX fkx_grp__grp_priv__grp_member ON grp_priv__grp_member(grp__id);
CREATE INDEX fkx_grp_priv__grp_priv__grp_member ON grp_priv__grp_member(grp_priv__id);

-- Everything below is left as a note for the future - in case we ever decided
-- actually allow privileges granted to individual users and/or individual
-- objects.

/*
-- 
-- TABLE: grp_priv__grp 
--

CREATE TABLE grp_priv__grp(
    grp_priv__id    INTEGER           NOT NULL,
    grp__id         INTEGER           NOT NULL,
    CONSTRAINT pk_grp_priv__grp PRIMARY KEY (grp_priv__id,grp__id)
) 
;


-- 
-- TABLE: grp_priv__person 
--

CREATE TABLE grp_priv__person(
    grp_priv__id    INTEGER           NOT NULL,
    person__id      INTEGER           NOT NULL,
    CONSTRAINT pk_grp_priv__person PRIMARY KEY (grp_priv__id,person__id)
) 
;


-- 
-- TABLE: grp_priv__usr 
--

CREATE TABLE grp_priv__usr(
    grp_priv__id    INTEGER           NOT NULL,
    usr__id        INTEGER           NOT NULL,
    CONSTRAINT pk_grp_priv__usr PRIMARY KEY (grp_priv__id,usr__id)
) 
;


-- 
-- TABLE: priv_table 
--

CREATE TABLE priv_table(
    id      INTEGER           NOT NULL,
    name    VARCHAR(30)    NOT NULL,
    CONSTRAINT pk_priv_table__id PRIMARY KEY (id)
) 
;


-- 
-- TABLE: usr_priv 
--

CREATE TABLE usr_priv(
    id          INTEGER           NOT NULL,
    usr__id    INTEGER           NOT NULL,
    value       INT2     NOT NULL,
    CONSTRAINT pk_usr_priv__id PRIMARY KEY (id)
) 
;


-- 
-- TABLE: usr_priv__grp 
--

CREATE TABLE usr_priv__grp(
    priv_usr__id    INTEGER           NOT NULL,
    grp__id          INTEGER           NOT NULL,
    CONSTRAINT pk_usr_priv__grp PRIMARY KEY (priv_usr__id,grp__id)
) 
;


-- 
-- TABLE: usr_priv__person 
--

CREATE TABLE usr_priv__person(
    usr_priv__id    INTEGER           NOT NULL,
    person__id       INTEGER           NOT NULL,
    CONSTRAINT pk_usr_priv__person PRIMARY KEY (usr_priv__id,person__id)
) 
;


-- 
-- TABLE: usr_priv__usr 
--

CREATE TABLE usr_priv__usr(
    usr_priv__id    INTEGER           NOT NULL,
    usr__id         INTEGER           NOT NULL,
    CONSTRAINT pk_usr_priv__usr PRIMARY KEY (usr_priv__id,usr__id)
) 
;


-- 
-- TABLE: usr_priv__grp_member 
--

CREATE TABLE usr_priv__grp_member(
    usr_priv__id    INTEGER           NOT NULL,
    grp__id          INTEGER           NOT NULL,
    CONSTRAINT pk_usr_priv__grp_member PRIMARY KEY (usr_priv__id,grp__id)
) 
;

*/


