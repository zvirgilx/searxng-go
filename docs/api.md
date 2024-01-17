### Search Api Definitions
#### Search from query

<details>
 <summary><code>GET</code> <code><b>/search</b></code><code>(get search result from query)</code></summary>

##### Parameters

> | name        | type     | data type | description                                              |
> |-------------|----------|-----------|----------------------------------------------------------|
> | q           | required | string    | query                                                    |
> | time_range  | option   | string    | time range of search result, e.g. day, week, mouth, year |
> | safe_search | option   | int       | search result content level                              |
> | language    | option   | string    | language, e.g. zh-CN, en-US, en-UK.                      |
> | category    | option   | string    | search category, e.g. general(default), video.           |
> | page_no     | option   | int       | the number of page, e.g. 1, 2, 3, ...                    |


##### Responses

> | name         | type         | data type       | description                   |
> |--------------|--------------|-----------------|-------------------------------|
> | query        | required     | string          | query                         |
> | results      | required     | list(Result)    | list of result                |
> | suggestions  | option(temp) | list(String)    | list of query suggestion      |
> | info_box     | option(temp) | object(InfoBox) | A information about the query |
> | next_page_no | required     | int             | next page_no of search page   |

Result

> | name      | type     | data type | description                           |
> |-----------|----------|-----------|---------------------------------------|
> | engine    | required | string    | engine name                           |
> | title     | required | string    | title                                 |
> | content   | required | string    | content                               |
> | url       | required | string    | url links to the third party          |
> | img_src   | option   | string    | image from result, e.g., movie poster |
> | thumbnail | option   | string    | thumbnail of video search result      |

InfoBox

> | name     | type     | data type  | description                                                         |
> |----------|----------|------------|---------------------------------------------------------------------|
> | title    | required | string     | title                                                               |
> | content  | required | string     | content                                                             |
> | url      | required | string     | url links to the detail of information, always is a third party url |
> | img_src  | option   | string     | image from result                                                   |
> | url_list | required | list(json) | url list to the third party.                                        |


##### ErrorCode

> | http code | content-type       | response                                    |
> |-----------|--------------------|---------------------------------------------|
> | `400`     | `application/json` | `{"code":"400","message":"Bad Request"}`    |
> | `500`     | `application/json` | `{"code":"500","message":"Internal Error"}` |

##### Example cURL

> ```javascript
>  curl -X GET 'http://localhost:8888/search?q=hello'
> ```

</details>

------------------------------------------------------------------------------------------
#### Auto query complete

<details>
 <summary><code>GET</code> <code><b>/complete</b></code><code>(complete the search query)</code></summary>

##### Parameters

> | name        | type     | data type | description                                              |
> |-------------|----------|-----------|----------------------------------------------------------|
> | q           | required | string    | query                                                    |

##### Responses

The completed language first depends on the `Accept Language` field in the browser request header, if not, the default language is used.

> | name    | type     | data type            | description                  |
> |---------|----------|----------------------|------------------------------|
> | query   | required | string               | query                        |
> | results | required | List(CompleteResult) | the complete result of query |

CompleteResult

> | name | type     | data type | description                             |
> |------|----------|-----------|-----------------------------------------|
> | type | required | string    | complete type, such as text, media,etc. |
> | text | required | string    | complete text from query                |
> | info | option   | string    | extra information                       |

##### ErrorCode

> | http code | content-type       | response                                    |
> |-----------|--------------------|---------------------------------------------|
> | `400`     | `application/json` | `{"code":"400","message":"Bad Request"}`    |
> | `500`     | `application/json` | `{"code":"500","message":"Internal Error"}` |

##### Example cURL

> ```javascript
>  curl -X GET 'http://localhost:8888/complete?q=hello'
> ```

</details>

------------------------------------------------------------------------------------------
