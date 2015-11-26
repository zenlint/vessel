## Vessel REST API Spec

### V1 Version

#### Workspace

##### **POST** `/v1/workspace`

##### **PUT** `/v1/workspace/:workspace`

##### **GET** `/v1/workspace/:workspace`

##### **DELETE** `/v1/workspace/:workspace`

#### Project

##### **POST** `/v1/project/:workspace`

##### **PUT** `/v1/project/:workspace/:project`

##### **GET** `/v1/project/:worksapce/:project`

##### **DELETE** `/v1/project/:worksapce/:project`

#### Pipeline

##### **POST** `/v1/pipeline/:project`

##### **PUT** `/v1/pipeline/:project/:pipeline`

##### **GET** `/v1/pipeline/:project/:pipeline`

##### **DELETE** `/v1/pipeline/:project/:pipeline`

##### **RUN** `/v1/pipeline/:project/:pipeline/run`

##### **GET** `/v1/pipeline/:project/:pipeline/status`

#### Point

##### **POST** `/v1/point/:pipeline`

##### **PUT** `/v1/point/:pipeline/:point`

##### **GET** `/v1/point/:pipeline/:point`

##### **DELETE** `/v1/point/:pipeline/:point`

#### Stage

##### **POST** `/v1/stage/:pipeline`

##### **PUT** `/v1/stage/:pipeline/:stage`

##### **GET** `/v1/stage/:pipeline/:stage`

##### **DELETE** `/v1/stage/:pipeline/:stage`

#### Param

##### **POST** `/v1/param/:uuid`

##### **GET** `/v1/param/:uuid/list`

##### **PUT** `/v1/param/:uuid/:param`

##### **GET** `/v1/param/:uuid/:param`

##### **DELETE** `/v1/param/:uuid/:param`
