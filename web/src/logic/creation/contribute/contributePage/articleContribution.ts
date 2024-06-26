import { getFullPathOfImage, getuploadingMethod } from "@/apis/commonality";
import { createArticleContribution, getArticleClassificationList, updateArticleContribution } from '@/apis/contribution';
import globalScss from "@/assets/styles/global/export.module.scss";
import { useEditArticleStore } from '@/store/creation';
import { GetUploadingMethodReq, GetUploadingMethodRes } from "@/types/commonality/commonality";
import { AreateArticleContributionReq, ArticleContribution, GetArticleClassificationListRes, GetArticleClassificationListResItem, UpdateArticleContributionReq } from '@/types/creation/contribute/contributePage/articleContribution';
import { uploadFileformation } from "@/types/creation/contribute/contributePage/vdeoContribution";
import { uploadFile } from '@/utils/upload/upload';
import { validateArticleTitle } from '@/utils/validate/validate';
import { ElInput, FormInstance, UploadProps, UploadRawFile, UploadRequestOptions } from 'element-plus';
import hljs from 'highlight.js';
import { QuillOptionsStatic } from 'quill';
import Swal from 'sweetalert2';
import { Ref, UnwrapNestedRefs, nextTick, reactive, ref, toRaw } from "vue";
import { Router, useRouter } from 'vue-router';

export const useArticleContributionProp = () => {
    const editArticleStore = useEditArticleStore()
    const router = useRouter()
    //Rich text component ref
    const myQuillEditor = ref()
    //content
    const form = reactive(<ArticleContribution>{
        id: 0,
        isShow: true,
        title: "",
        content: "",
        comments: true,
        labelInputVisible: false,
        labelText: "",
        label: [],
        classification_id: 0,
    })

    const labelInputRef = ref<InstanceType<typeof ElInput>>()
    //Configuration
    const options = reactive(<QuillOptionsStatic>{
        modules: {
            toolbar: {
                container: [
                    [{ size: ["small", false, "large"] }],
                    ["bold", "italic", "underline"],
                    ['code-block'],
                    [{ header: 1 }, { header: 2 }],
                    [{ list: "ordered" }, { list: "bullet" }],
                    ["link", "image"],
                    [{ color: [] }, { background: [] }],
                    [{ align: [] }]
                ],
                handlers: {
                    image: function (value: any) {
                        console.log("value", value)
                        if (value) {
                            // Trigger a custom upload
                            console.log("Trigger a custom upload")
                            console.log(uploadBtnRef.value.$el.click())
                        } else {
                            myQuillEditor.value.format("image", false);
                        }
                    }
                },
                image: () => {

                }
            },
            history: {
                delay: 1000,
                maxStack: 50,
                userOnly: false
            },
            syntax: {
                highlight: (text: string) => hljs.highlightAuto(text).value
            },
        },
    });
    //Upload component ref
    const uploadRef = ref()
    const uploadBtnRef = ref()
    const ruleFormRef = ref<FormInstance>()
    const uploadProgressRef = ref<HTMLElement | null>()
    //Rich text image upload
    const uploadFileformation = reactive(<uploadFileformation>{
        progress: 0,
        FileUrl: '',
        uploadUrl: "",
        interface: "articleContribution",
        uploadType: "",
        action: "#",
    })
    //Article cover upload
    const uploadCoveration = reactive(<uploadFileformation>{
        progress: 0,
        FileUrl: '',
        uploadUrl: "",
        interface: "articleContributionCover",
        uploadType: "",
        action: "#",
    })

    const colors = [
        { color: '#f56c6c', percentage: 20 },
        { color: '#e6a23c', percentage: 40 },
        { color: '#5cb87a', percentage: 60 },
        { color: '#1989fa', percentage: 80 },
        { color: '#6f7ad3', percentage: 100 },
    ];
    const cascaderValue = ref([])
    const cascader = ref(<GetArticleClassificationListRes>[])

    return {
        router,
        uploadFileformation,
        uploadCoveration,
        labelInputRef,
        ruleFormRef,
        uploadRef,
        uploadBtnRef,
        uploadProgressRef,
        myQuillEditor,
        form,
        options,
        colors,
        cascaderValue,
        cascader,
        editArticleStore
    }
}

//Upload file processing
export const useHandleFileMethod = (uploadFileformation: uploadFileformation, form: ArticleContribution, myQuillEditor: Ref, uploadProgressRef: Ref) => {

    const handleFileSuccess: UploadProps['onSuccess'] = async (
        response,
        uploadFile
    ) => {
        uploadFileformation.FileUrl = URL.createObjectURL(uploadFile.raw!)
        //Splice full path
        const path = await getFullPathOfImage({ path: uploadFileformation.uploadUrl, type: uploadFileformation.uploadType })
        // Get rich text instance
        let quill = toRaw(myQuillEditor.value).getQuill()
        // Get cursor position
        let length = quill.selection.savedRange.index;
        // Insert picture 
        quill.insertEmbed(length, "image", path.data);
        // Adjust the cursor to the end 
        quill.setSelection(length + 1);

        Swal.close()
    }

    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {
        console.log("upload failed")
        Swal.fire({
            title: "upload failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",

        })
        console.log(response)

    }

    //Processing before uploading
    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile: UploadRawFile) => {
        return true
    }

    //Modify default request
    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            console.log(uploadProgressRef.value.$el)
            uploadProgressRef.value.$el.style.display = "";
            Swal.fire({
                title: 'File uploading...',
                html: uploadProgressRef.value.$el,
                showClass: {
                    popup: 'animate__animated animate__fadeInDown'
                },
                hideClass: {
                    popup: 'animate__animated animate__fadeOutUp'
                },
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
            })
            const response = await uploadFile(uploadFileformation, params.file)

            console.log("Upload successful")
            console.log(response)
            uploadFileformation.uploadUrl = response.path
        } catch (err) {
            console.log(err)
            Swal.fire({
                title: "Failed to upload image",
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
//Upload cover processing
export const useHandleCoverMethod = (uploadCoveration: uploadFileformation) => {

    const handleFileSuccess: UploadProps['onSuccess'] = (
        response,
        uploadFile
    ) => {
        uploadCoveration.FileUrl = URL.createObjectURL(uploadFile.raw!)
    }

    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {
        console.log("upload failed")
        Swal.fire({
            title: "upload failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",

        })
        console.log(response)

    }

    //Processing before uploading
    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile: UploadRawFile) => {
        return await new Promise<boolean>((resolve, reject) => {
            //Determine size
            if (rawFile.size / 1024 / 1024 > 2) {
                Swal.fire({
                    title: "Cover size cannot be larger than 2 m",
                    heightAuto: false,
                    icon: "error",

                })
                reject(false);
            }
            //Determine size
            let reader = new FileReader();
            reader.readAsDataURL(rawFile);
            reader.onload = () => {
                // Let the src of the img tag in the page point to the read path
                let img = new Image();
                img.src = reader.result as string;
                img.onload = () => {
                    console.log(img.width);
                    console.log(img.height);
                    if (img.width < 960 || img.height < 600) {
                        Swal.fire({
                            title: "Please upload images above 960*600 size",
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

    //Modify default request
    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            const response = await uploadFile(uploadCoveration, params.file)
            console.log(response)
            uploadCoveration.uploadUrl = response.path
        } catch (err) {
            console.log(err)
            Swal.fire({
                title: "Failed to upload cover",
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

//Tag handling
export const userLabelHandlMethod = (form: ArticleContribution, labelInputRef: Ref) => {
    const handleClose = (tag: string) => {
        form.label.splice(form.label.indexOf(tag), 1)
    }

    const showInput = () => {
        form.labelInputVisible = true
        nextTick(() => {
            labelInputRef.value!.input!.focus()
        })
    }

    const handleInputConfirm = () => {
        if (form.labelText) {
            form.label.push(form.labelText)
        }
        form.labelInputVisible = false
        form.labelText = ''
    }

    return {
        handleClose,
        showInput,
        handleInputConfirm
    }

}
export const useSaveData = async (form: ArticleContribution, formEl: FormInstance | undefined, router: Router, uploadFileformation: uploadFileformation, uploadCoveration: uploadFileformation, props: any) => {
    if (!formEl) return;
    await formEl.validate(async (valid, fields) => {
        if (valid) {
            try {
                if (!uploadCoveration.uploadUrl) {
                    Swal.fire({
                        title: "Please upload the cover first",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "warning",
                    })
                    return
                }
                console.log(form.classification_id, form.classification_id == 0)
                if (form.classification_id === 0) {
                    Swal.fire({
                        title: "please select a type",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "warning",
                    })
                    return
                }
                if (props.type == "edit") {
                    //update mode
                    let updateRequistData = <UpdateArticleContributionReq>{
                        id: form.id,
                        cover: uploadCoveration.uploadUrl,
                        coverUploadType: uploadCoveration.uploadType,
                        articleContributionUploadType: uploadFileformation.uploadType,
                        title: form.title,
                        label: form.label,
                        content: form.content,
                        comments: form.comments,
                        classification_id: form.classification_id
                    }
                    await updateArticleContribution(updateRequistData)

                } else {
                    let createRequistData = <AreateArticleContributionReq>{
                        cover: uploadCoveration.uploadUrl,
                        coverUploadType: uploadCoveration.uploadType,
                        articleContributionUploadType: uploadFileformation.uploadType,
                        title: form.title,
                        label: form.label,
                        content: form.content,
                        comments: form.comments,
                        classification_id: form.classification_id
                    }
                    await createArticleContribution(createRequistData)
                }
                let swalTitle = props.type == "edit" ? "update completed" : "Posted successfully"
                Swal.fire({
                    title: swalTitle,
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "success",
                    preConfirm: () => {
                        router.push({ name: "ArticleManagement" })
                    }
                })
            } catch (err) {
                console.log(err)
                Swal.fire({
                    title: (err as Error).message,
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




export const useInit = async (uploadFileformation: uploadFileformation, uploadCoveration: uploadFileformation, cascader: Ref<GetArticleClassificationListRes>, form: UnwrapNestedRefs<ArticleContribution>, editArticleStore: any, props: any, myQuillEditor: Ref<any>, cascaderValue: Ref<Array<number>>) => {
    try {
        //Get the request method of the current interface
        const updataMenhod = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: uploadFileformation.interface
        })).data as GetUploadingMethodRes
        uploadFileformation.uploadType = updataMenhod.type
        const updataMenhodCover = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: uploadCoveration.interface
        })).data as GetUploadingMethodRes
        uploadCoveration.uploadType = updataMenhod.type
        //Get article classification
        const cn = await getArticleClassificationList()
        if (cn.data?.length == 0) {
            Swal.fire({
                title: "Get the article category is empty",
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
                icon: "error",
            })
            return false
        }
        cascader.value = cn.data as GetArticleClassificationListRes

        if (props.type == "edit") {
            //edit mode
            form.id = editArticleStore.editArticleData.articleID
            form.title = editArticleStore.editArticleData.title
            myQuillEditor.value.setContents(editArticleStore.editArticleData.content)
            form.comments = editArticleStore.editArticleData.is_comments
            form.label = editArticleStore.editArticleData.label
            form.classification_id = editArticleStore.editArticleData.classification_id
            uploadCoveration.uploadType = editArticleStore.editArticleData.cover_upload_type
            uploadCoveration.uploadUrl = editArticleStore.editArticleData.cover_url
            uploadCoveration.FileUrl = editArticleStore.editArticleData.cover
            //Convert multidimensional array to one-dimensional array
            let oneDimensionalArr: Array<GetArticleClassificationListResItem> = []
            classificationSuperior(oneDimensional(cascader.value, oneDimensionalArr), cascaderValue.value, form.classification_id)
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
    const articleContributionRules = reactive({
        title: [{ validator: validateArticleTitle, trigger: 'change' }],
    });
    return {
        articleContributionRules,
    };
};



//Find the superior node for distribution
function classificationSuperior(data: GetArticleClassificationListRes, arr: Array<number>, aid: number) {
    if (aid == 0) return arr
    let node = data.filter((item) => {
        return item.id == aid
    })[0]
    arr.unshift(node.id)
    classificationSuperior(data, arr, node.aid)
}

function oneDimensional(data: Array<GetArticleClassificationListResItem>, arr: Array<GetArticleClassificationListResItem>) {
    data.filter((item: GetArticleClassificationListResItem) => {
        arr.push(item)
        if (item.children == null) return
        if (item.children.length > 0) {
            oneDimensional(item.children, arr)
        }
    })
    return arr
}