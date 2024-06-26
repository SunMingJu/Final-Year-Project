import { deleteVideoByID, getVideoManagementList } from '@/apis/contribution';
import globalScss from "@/assets/styles/global/export.module.scss";
import { useEditVideoStore } from '@/store/creation';
import { useGlobalStore } from '@/store/main';
import { DeleteVideoByIDReq, GetVideoManagementListItem, GetVideoManagementListRes } from '@/types/creation/manuscript/video';
import { PageInfo } from '@/types/idnex';
import { GetRecordListReq } from '@/types/personal/record/record';
import { editVideo } from '@/types/store/creation';
import Swal from 'sweetalert2';
import { Ref, ref } from 'vue';
import { Router, useRoute, useRouter } from 'vue-router';

export const useVideoProp = () => {
    const loading = useGlobalStore().globalData.loading
    const editVideoStore = useEditVideoStore()
    const route = useRoute()
    const router = useRouter()
    const videoList = ref(<GetVideoManagementListRes>[])
    const pageInfo = ref(<PageInfo>{
        page: 1,
        size: 8,
    })
 //Whether it is loaded for the first time
    const isLoading = ref(true)
    //Is loading more
    const isLoadMore = ref(false)
    //Whether all loading is completed
    const isTheEnd = ref(false)
    return {
        videoList,
        isLoading,
        route,
        router,
        pageInfo,
        loading,
        isLoadMore,
        isTheEnd,
        editVideoStore
    }
}

export const useJump = (item: GetVideoManagementListItem, router: Router) => {
    router.push({ name: "VideoShow", params: { id: item.id } })
}

export const useDelRecord = async (recordList: Ref<GetVideoManagementListRes>, id: number) => {
    try {
        //Delete request
        Swal.fire({
            heightAuto: false,
            title: 'Confirm to delete this video',
            icon: 'warning',
            confirmButtonColor: globalScss.colorButtonTheme,
            showCancelButton: true,
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
            reverseButtons: true
        }).then(async (result) => {
            if (result.isConfirmed) {
                try {
                    await deleteVideoByID(<DeleteVideoByIDReq>{
                        id
                    })
                    Swal.fire({
                        title: "successfully deleted",
                        confirmButtonColor: globalScss.colorButtonTheme,
                        heightAuto: false,
                        icon: "success",
                    })
                    recordList.value = recordList.value.filter((item: GetVideoManagementListItem) => {
                        if (item.id == id) item.is_delete = true
                        return item
                    })
                    setTimeout(() => {
                        recordList.value = recordList.value.filter((item: GetVideoManagementListItem) => {
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


export const useEditRecord = async (item: GetVideoManagementListItem, loading: any, editVideoStore: any, router: Router) => {
    try {
        editVideoStore.setEditVideoData(<editVideo>{
            videoID: item.id,
            cover: item.cover,
            reprinted: item.reprinted,
            cover_url: item.cover_url,
            coverUploadType: item.cover_upload_type,
            title: item.title,
            label: item.label,
            introduce: item.introduce,
            videoDuration: item.video_duration
        })
        router.push({ name: "Contribute", query: { type: "editVideo" } })
    } catch (err) {
        Swal.fire({
            title: "Edit failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}


export const useLoadData = async (videoList: Ref<GetVideoManagementListRes>, isLoading: Ref<boolean>, page: Ref<PageInfo>, isTheEnd: Ref<boolean>,) => {
    try {
        const data = await getVideoManagementList(<GetRecordListReq>{
            page_info: page.value
        })
        page.value.page++
        if (!data.data) return false
        if (data.data.length == 0) isTheEnd.value = true
        data.data = data.data.filter((item) => {
            item.is_delete = false
            return item
        })
        videoList.value = [...videoList.value, ...data.data]
        console.log(videoList)
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