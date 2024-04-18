import { getuploadingMethod } from "@/apis/commonality";
import { createFavorites } from "@/apis/personal";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useUserStore } from "@/store/main";
import { GetUploadingMethodReq, GetUploadingMethodRes } from "@/types/commonality/commonality";
import { CreateCollectRmation, SaveCreateFavoritesDataReq } from "@/types/personal/collect/createFavorites";
import { GetFavoritesListItem } from "@/types/personal/collect/favorites";
import { GetLiveDataRes } from "@/types/personal/live/setUp";
import { uploadFile } from '@/utils/upload/upload';
import { validateCollectTitle } from "@/utils/validate/validate";
import { FormInstance, UploadProps, UploadRequestOptions } from 'element-plus';
import Swal from 'sweetalert2';
import { reactive, ref } from "vue";
import { Router, useRouter } from 'vue-router';

export const useCreateFavoritesProp = () => {
    const userStore = useUserStore()
    const router = useRouter()
    const saveDateFormRef = ref<FormInstance>()
    const createFavoriteRmationForm = reactive<CreateCollectRmation>({
        id: 0,
        FileUrl: '',
        uploadUrl: "",
        interface: "createFavoritesCover",
        title: "",
        content: "",
        uploadType: "",
        action: "#",
        progress: 0
    });

    //Define request result raw data
    const rawData = reactive<GetLiveDataRes>({
        title: "",
        img: ""
    })

    return {
        userStore,
        createFavoriteRmationForm,
        saveDateFormRef,
        rawData,
        router
    }
}


//File upload processing
export const useHandleFileMethod = (createFavoriteRmationForm: CreateCollectRmation) => {
    const handleFileSuccess: UploadProps['onSuccess'] = (
        response,
        uploadFile
    ) => {
        createFavoriteRmationForm.FileUrl = URL.createObjectURL(uploadFile.raw!)
    }

    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {console.log("Upload failed")
        Swal.fire({
            title: "Upload failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
        console.log(response)
    }


    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile) => {
        return await new Promise<boolean>((resolve, reject) => {
            //Judge the size
            if (rawFile.size /1024 /1024 > 2) {
                Swal.fire({
                    title: "Cover size cannot be larger than 2M",
                    heightAuto: false,
                    icon: "error",

                })
                reject(false);
            }
            //Judge size
            let reader = new FileReader();
            reader.readAsDataURL(rawFile);
            reader.onload = () => {
                //Let the src of the img tag in the page point to the read path
                let img = new Image();
img.src = reader.result as string;
                img.onload = () => {
                    console.log(img.width);
                    console.log(img.height);
                    if (img.width < 960 || img.height < 600) {
                        Swal.fire({
                            title: "Please upload pictures above 960*600 size",
                            heightAuto: false,
                            confirmButtonColor: globalScss.colorButtonTheme,
                            icon: "error",
                        });
                        reject(false);
                    } else {
                        resolve(true);
                    }
                };
            };
        })
    }
    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            const response = await uploadFile(createFavoriteRmationForm, params.file)
            console.log(response)
            createFavoriteRmationForm.uploadUrl = response?.path
            console.log(response)
        } catch (err) {
            console.log(err)
            Swal.fire({
                title: "Failed to obtain upload node",
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
                icon: "error",
            })
        }
    }

    return {
        handleFileSuccess,
        beforeFileUpload,
        handleFileError,
        RedefineUploadFile
    }

}

export const useSaveData = async (createFavoriteRmationForm: CreateCollectRmation, formEl: FormInstance | undefined, rawData: GetLiveDataRes, router: Router, emits: any) => {
    if (!formEl) return;

    await formEl.validate(async (valid, fields) => {
        if (valid) {
            try {
                if (createFavoriteRmationForm.uploadUrl == rawData.img && createFavoriteRmationForm.title == rawData.title) throw "Unmodified information";
                if (!createFavoriteRmationForm.uploadUrl) throw "Please upload pictures first"
                const requistData = <SaveCreateFavoritesDataReq>{
                    id: createFavoriteRmationForm.id,
                    type: createFavoriteRmationForm.uploadType,
                    cover: createFavoriteRmationForm.uploadUrl,
                    title: createFavoriteRmationForm.title,
                    content: createFavoriteRmationForm.content
                }
                await createFavorites(requistData)
                emits("shutDown");
                if (createFavoriteRmationForm.id == 0) {
                    //clear content
                    createFavoriteRmationForm.FileUrl = ""
                    createFavoriteRmationForm.title = ""
                    createFavoriteRmationForm.content = ""
                    //Create mode
                    Swal.fire({
                        title: "Created successfully",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "success",
preConfirm: () => {
                            router.push({ name: "Favorites", query: { type: 'createTime' + Date.now() } })
                        }
                    })
                } else {
                    //Update model
                    Swal.fire({
                        title: "Update Success",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "success",
                    })
                }

            } catch (err: any) {
                console.log(err)
                emits("shutDown");
                Swal.fire({
                    title: err.message as string,
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "warning",
                })
            }
        } else {
            console.log('error submit!', fields)
        }
    })
}


export const useInit = async (createFavoriteRmationForm: CreateCollectRmation, rawData: GetLiveDataRes, type: boolean, item: GetFavoritesListItem | undefined) => {
    try {
        //Get the request method of the current interface
        const updataMenhod = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: createFavoriteRmationForm.interface
        })).data as GetUploadingMethodRes
        createFavoriteRmationForm.uploadType = updataMenhod.type
        if (!type) {
            if (item == undefined) return false
            if (createFavoriteRmationForm.FileUrl) {
                //Change the upload type only after uploading the image
                createFavoriteRmationForm.uploadType = item.type
            }
            createFavoriteRmationForm.title = item.title
            createFavoriteRmationForm.content = item.content
            createFavoriteRmationForm.uploadUrl = item.src
            createFavoriteRmationForm.FileUrl = item.cover
            createFavoriteRmationForm.id = item.id
        }

    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to get upload method",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}


//form validation
export const useRules = () => {
    const liveInformationRules = reactive({
        title: [{ validator: validateCollectTitle, trigger: 'change' }],
    });
    return {
        liveInformationRules,
    };
};