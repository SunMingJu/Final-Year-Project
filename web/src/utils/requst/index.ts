import router from "@/router";
import { useUserStore } from '@/store/main';
import { FileSliceUpload, FileUpload } from '@/types/idnex';
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { ElMessage } from 'element-plus';
import Swal from "sweetalert2";
//Interface for data return
//Define request response parameters, excluding data
const Toast = Swal.mixin({
  toast: true,
  position: 'top',
  showConfirmButton: false,
  timer: 3000,
})
interface Result {
  code: number;
  message: string
}

//Request response parameters, including datainterface ResultData<T = any> extends Result {
  data?: T;
}


const URL: string = import.meta.env.VITE_BASE_URL;

enum RequestEnums {
  TIMEOUT = 60000,
  SUCCESS = 200, //Request successful
  OPERATIONFAIL = 500, //Operation failed
  NOTLOGIN = 303, //Operation failed
  FAIL = 999, //Request failed
}
const config = {
  //default address
  baseURL: URL as string,
  //Set the timeout
  timeout: RequestEnums.TIMEOUT as number,

}

const userInfo = useUserStore()

class RequestHttp {
  // Define member variables and specify types
  service: AxiosInstance;
  public constructor(config: AxiosRequestConfig, user: any) {

    // Instantiate axios
    this.service = axios.create(config);

    /**
     *Request interceptor
     *Client sends request -> [Request Interceptor] -> Server
     *Token verification (JWT): Accept the token returned by the server and store it in vuex/pinia/local storage
     */
    this.service.interceptors.request.use(
      (config: AxiosRequestConfig) => {
        const token = user.userInfoData.token || '';
        return {
          ...config,
          headers: {
            'token': token, //The request header carries token information
          }
        }
      },
      (error: AxiosError) => {
        // Request error report
        Promise.reject(error)
      }
    )

   /**
     *Response interceptor
     *The server returns information -> [Unified interception processing] -> Client JS obtains the information
     */
    this.service.interceptors.response.use(
      (response: AxiosResponse) => {
        const { data, config } = response; // deconstruct
        if (data.code == RequestEnums.OPERATIONFAIL) {
          return Promise.reject(data);
        }
        if (data.code === RequestEnums.NOTLOGIN) {
          // If the login information is invalid, you should jump to the login page and clear the local token.
          userInfo.userInfoData.token = ""
          router.push({
            path: '/login'
          })
          return Promise.reject(data);
        }
        // Global error message interception (to prevent the data stream from being returned when downloading files, without code, and directly reporting errors)
        if (data.code && data.code !== RequestEnums.SUCCESS) {
          ElMessage.error(data); // You can also use components here to prompt error messages
          return Promise.reject(data)
        }
        return data;
      },
      (error: AxiosError) => {
        const { response } = error;
        if (response) {
          this.handleCode(response.status)
        }
        if (!window.navigator.onLine) {
          Toast.fire({
            icon: 'error',
            title: 'Network connection failed'
          })
          //You can jump to the error page or do nothing.
          //return router.replace({
          //path: '/404'
          //});
        }
      }
    )
  }

  handleCode(code: number): void {
    switch (code) {
      case 401:
        Toast.fire({
          icon: 'error',
          title: 'Login failed, please log in again'
        })
        break;
      default:
        Toast.fire({
          icon: 'error',
          title: 'Request failed'
        })
        break;
    }
  }

  // Common method encapsulation
  get<T>(url: string, params?: object): Promise<ResultData<T>> {
    return this.service.get(url, { params });
  }
  post<T>(url: string, params?: object): Promise<ResultData<T>> {
    return this.service.post(url, params);
  }
  upload<T>(url: string, params: object, uploadConfig: FileUpload): Promise<ResultData<T>> {
    return this.service.post(url, params, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: ProgressEvent => {
        if (!ProgressEvent?.total) return;
        //Calculate progress bar
        uploadConfig.progress = Math.round(ProgressEvent.loaded / ProgressEvent?.total * 100)
      }
    });
  }
  uploadSlice<T>(url: string, params: object, uploadConfig: FileSliceUpload): Promise<ResultData<T>> {
    return this.service.post(url, params, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: ProgressEvent => {
        if (!ProgressEvent?.total) return;
        console.log(Math.round(ProgressEvent.loaded / ProgressEvent?.total * 100))
        //Calculate progress bar
        uploadConfig.progress = Math.round(ProgressEvent.loaded / ProgressEvent?.total * 100)
      }
    });
  }
  put<T>(url: string, params?: object): Promise<ResultData<T>> {
    return this.service.put(url, params);
  }
  delete<T>(url: string, params?: object): Promise<ResultData<T>> {
    return this.service.delete(url, { params });
  }
}

// Export an instance object
export default new RequestHttp(config, userInfo);