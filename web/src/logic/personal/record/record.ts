import { deleteRecordByID, getRecordList } from '@/apis/personal';
import globalScss from "@/assets/styles/global/export.module.scss";
import { useGlobalStore } from '@/store/main';
import { PageInfo } from '@/types/idnex';
import { DeleteRecordByIDReq, GetRecordListItem, GetRecordListReq, GetRecordListRes } from '@/types/personal/record/record';
import Swal from 'sweetalert2';
import { Ref, ref } from 'vue';
import { Router, useRoute, useRouter } from 'vue-router';

export const useRecordProp = () => {
    const loading = useGlobalStore().globalData.loading
    const route = useRoute()
    const router = useRouter()
    const recordList = ref(<GetRecordListRes>[])
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
        recordList,
        isLoading,
        route,
        router,
        pageInfo,
        loading,
        isLoadMore,
        isTheEnd
    }
}

export const useJump = (item: GetRecordListItem, router: Router) => {
    if (item.type == "video") {
        router.push({ name: "VideoShow", params: { id: item.to_id } })
    } else if (item.type == "Column") {
        router.push({ name: "ArticleShow", query: { articleID: item.to_id } })
    } else {
        router.push({ name: "liveRoom", query: { roomID: item.to_id } })
    }
}

export const useDelRecord = async (recordList: Ref<GetRecordListRes>, id: number, loading: any) => {
    try {
        loading.loading = true
        await deleteRecordByID(<DeleteRecordByIDReq>{
            id
        })
        loading.loading = false
        recordList.value = recordList.value.filter((item) => {
            if (item.id == id) item.is_delete = true
            return item
        })
        setTimeout(() => {
            recordList.value = recordList.value.filter((item) => {
                return item.id != id
            })
        }, 400)

    } catch (err) {
        loading.loading = false
        Swal.fire({
            title: "failed to delete",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}


export const useLoadData = async (recordList: Ref<GetRecordListRes>, isLoading: Ref<boolean>, page: Ref<PageInfo>, isTheEnd: Ref<boolean>,) => {
    try {
        const data = await getRecordList(<GetRecordListReq>{
            page_info: page.value
        }
        )
        page.value.page++
        if (!data.data) return false
        if (data.data.length == 0) isTheEnd.value = true
        data.data = data.data.filter((item) => {
            item.is_delete = false
            return item
        })
        recordList.value = [...recordList.value, ...data.data]
        console.log(recordList)
        isLoading.value = false
    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to get history",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}