# media_types

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE "media_types"
(
    [MediaTypeId] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    [Name] NVARCHAR(120)
)
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| MediaTypeId | INTEGER |  | false | [tracks](tracks.md) |  |  |
| Name | NVARCHAR(120) |  | true |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| MediaTypeId | PRIMARY KEY | PRIMARY KEY (MediaTypeId) |

## Relations

![er](media_types.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
