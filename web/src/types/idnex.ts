//http
export interface Result {
  code: number;
  msg: string
}
// Request response parameters, including data
export interface ResultData<T = any> extends Result {
  data?: T;
}
export interface ResultWs {
  code: number;
  type: string;
  message: string;
}
// Request response parameters, including data
export interface ResultDataWs<T = any> extends ResultWs {
  data?: T;
}


// Oss configuration
export interface OssConfig {
  access_id: string;
  host: string;
  expire: number;
  signature: string;
  policy: string;
  dir: string;
}

export interface OssSTSInfo {
  region: string;
  accessKeyId: string;
  accessKeySecret: string;
  stsToken: string;
  bucket: string;
  expirationTime: number;
}


// File upload configuration is required
export interface FileUpload {
  progress: number; //Upload progress
  FileUrl: string; //return file path
  interface: string; //Upload interface name
  uploadUrl: string; //Upload path
  uploadType: string; //Upload type
  action: string; //Request address
}


//File upload configuration is required
export interface FileSliceUpload {
  index: number;
  progress: number; //Upload progress
  size: number;
}

//Paging configuration
export interface PageInfo {
  page: number; //page number
  size: number; //size of each page
  keyword?: string; //keyword
}