# tracks

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE "tracks"
(
    [TrackId] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    [Name] NVARCHAR(200)  NOT NULL,
    [AlbumId] INTEGER,
    [MediaTypeId] INTEGER  NOT NULL,
    [GenreId] INTEGER,
    [Composer] NVARCHAR(220),
    [Milliseconds] INTEGER  NOT NULL,
    [Bytes] INTEGER,
    [UnitPrice] NUMERIC(10,2)  NOT NULL,
    FOREIGN KEY ([AlbumId]) REFERENCES "albums" ([AlbumId]) 
		ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY ([GenreId]) REFERENCES "genres" ([GenreId]) 
		ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY ([MediaTypeId]) REFERENCES "media_types" ([MediaTypeId]) 
		ON DELETE NO ACTION ON UPDATE NO ACTION
)
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| TrackId | INTEGER |  | false | [invoice_items](invoice_items.md) [playlist_track](playlist_track.md) |  |  |
| Name | NVARCHAR(200) |  | false |  |  |  |
| AlbumId | INTEGER |  | true |  | [albums](albums.md) |  |
| MediaTypeId | INTEGER |  | false |  | [media_types](media_types.md) |  |
| GenreId | INTEGER |  | true |  | [genres](genres.md) |  |
| Composer | NVARCHAR(220) |  | true |  |  |  |
| Milliseconds | INTEGER |  | false |  |  |  |
| Bytes | INTEGER |  | true |  |  |  |
| UnitPrice | NUMERIC(10,2) |  | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| TrackId | PRIMARY KEY | PRIMARY KEY (TrackId) |
| - (Foreign key ID: 0) | FOREIGN KEY | FOREIGN KEY (MediaTypeId) REFERENCES media_types (MediaTypeId) ON UPDATE NO ACTION ON DELETE NO ACTION MATCH NONE |
| - (Foreign key ID: 1) | FOREIGN KEY | FOREIGN KEY (GenreId) REFERENCES genres (GenreId) ON UPDATE NO ACTION ON DELETE NO ACTION MATCH NONE |
| - (Foreign key ID: 2) | FOREIGN KEY | FOREIGN KEY (AlbumId) REFERENCES albums (AlbumId) ON UPDATE NO ACTION ON DELETE NO ACTION MATCH NONE |

## Indexes

| Name | Definition |
| ---- | ---------- |
| IFK_TrackMediaTypeId | CREATE INDEX [IFK_TrackMediaTypeId] ON "tracks" ([MediaTypeId]) |
| IFK_TrackGenreId | CREATE INDEX [IFK_TrackGenreId] ON "tracks" ([GenreId]) |
| IFK_TrackAlbumId | CREATE INDEX [IFK_TrackAlbumId] ON "tracks" ([AlbumId]) |

## Relations

![er](tracks.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)