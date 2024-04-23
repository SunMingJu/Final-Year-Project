import { deleteArticleByID, getArticleManagementList } from '@/apis/contribution';
import globalScss from "@/assets/styles/global/export.module.scss";
import { useEditArticleStore } from '@/store/creation';
import { useGlobalStore } from '@/store/main';
import { DeleteArticleByIDReq, GetArticleManagementListItem, GetArticleManagementListRes } from '@/types/creation/manuscript/Article';
import { PageInfo } from '@/types/idnex';
import { GetRecordListReq } from '@/types/personal/record/record';
import { editArticle } from '@/types/store/creation';
import Swal from 'sweetalert2';
import { Ref, ref } from 'vue';
import { Router, useRoute, useRouter } from 'vue-router';

export const useArticleProp = () => {
    const loading = useGlobalStore().globalData.loading
    const editArticleStore = useEditArticleStore()
    const route = useRoute()
    const router = useRouter()
    const articleList = ref(<GetArticleManagementListRes>[])
    const pageInfo = ref(<PageInfo>{
        page: 1,
        size: 8,
    })
    //Whether to load for the first time
    const isLoading = ref(true)
    //Is loading more
    const isLoadMore = ref(false)
    //Whether all loading is completed
    const isTheEnd = ref(false)
    return {
        articleList,
        isLoading,
        route,
        router,
        pageInfo,
        loading,
        isLoadMore,
        isTheEnd,
        editArticleStore
    }
}

export const useJump = (item: GetArticleManagementListItem, router: Router) => {
    router.push({ name: "ArticleShow", params: { id: item.id } })
}

export const useDelRecord = async (recordList: Ref<GetArticleManagementListRes>, id: number) => {
    try {
        //Delete request
        Swal.fire({
            heightAuto: false,
            title: 'Are you sure you want to delete this article?',
            icon: 'warning',
            confirmButtonColor: globalScss.colorButtonTheme,
            showCancelButton: true,
            confirmButtonText: 'confirm',
            cancelButtonText: 'Cancel',
            reverseButtons: true
        }).then(async (result) => {
            if (result.isConfirmed) {
                try {
                    await deleteArticleByID(<DeleteArticleByIDReq>{
                        id
                    })
                    Swal.fire({
                        title: "successfully deleted",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "success",
                    })
                    recordList.value = recordList.value.filter((item) => {
                        if (item.id == id) item.is_delete = true
                        return item
                    })
                    setTimeout(() => {
                        recordList.value = recordList.value.filter((item) => {
                            return item.id != id
                        })
                    }, 400)
                } catch (err: any) {
                    console.log(err)
                    Swal.fire({
                        title: "failed to delete",
                        heightAuto: false,
                        confirmButtonColor: globalScss.colorButtonTheme,
                        icon: "error",
                    })
                }
            } else if (result.dismiss === Swal.DismissReason.cancel) {
                console.log("Cancel")
            }
        })

    } catch (err) {

        Swal.fire({
            title: "failed to delete",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}


export const useEditRecord = async (item: GetArticleManagementListItem, loading: any, editArticleStore: any, router: Router) => {
    try {
        editArticleStore.setEditArticleData(<editArticle>{
            articleID: item.id,
            title: item.title,
            cover: item.cover,
            cover_url: item.cover_url,
            cover_upload_type: item.cover_upload_type,
            content: item.content,
            is_comments: item.is_comments,
            label: item.label,
            classification_id: item.classification_id,
        })
        router.push({ name: "Contribute", query: { type: "editArticle" } })
    } catch (err) {
        Swal.fire({
            title: "Edit failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}


export const useLoadData = async (articleList: Ref<GetArticleManagementListRes>, isLoading: Ref<boolean>, page: Ref<PageInfo>, isTheEnd: Ref<boolean>,) => {
    try {
        const data = await getArticleManagementList(<GetRecordListReq>{
            page_info: page.value
        })
        page.value.page++
        if (!data.data) return false
        if (data.data.length == 0) isTheEnd.value = true
        data.data = data.data.filter((item) => {
            item.is_delete = false
            return item
        })
        articleList.value = [...articleList.value, ...data.data]
        console.log(articleList)
        isLoading.value = false
    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to obtain video manuscript",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}