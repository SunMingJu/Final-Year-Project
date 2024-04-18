import { deleteFavorites, getFavoritesList } from "@/apis/personal";
import { useUserStore } from "@/store/main";
import Swal from "sweetalert2";
import { useRoute, useRouter } from "vue-router";
import globalScss from "@/assets/styles/global/export.module.scss";
import { Ref, ref } from "vue";
import { DeleteFavoritesReq, GetFavoritesListItem, GetFavoritesListRes } from "@/types/personal/collect/favorites";


export const useFavoritesProp = () => {
    const userStore = useUserStore()
    const router = useRouter()
    const route = useRoute()
    const createCollectDialogShow = ref(false)
    const favoritesList = ref(<GetFavoritesListRes>{})
    const updataInfo = ref(<GetFavoritesListItem>{})
    return {
        userStore,
        router,
        favoritesList,
        createCollectDialogShow,
        updataInfo,
        route,
    }
}


export const useDelFavorites = (id: number , favoritesList : Ref<GetFavoritesListRes>) => {
    Swal.fire({
        heightAuto: false,
        title: 'Are you sure you want to delete this favorite?',
        icon: 'warning',
        confirmButtonColor: globalScss.colorButtonTheme,
        showCancelButton: true,
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        reverseButtons: true
    }).then(async (result) => {
        if (result.isConfirmed) {
            try {
                await deleteFavorites(<DeleteFavoritesReq>{
                    id
                })
                Swal.fire({
                    title: "successfully deleted",
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "success",
                })
                favoritesList.value = favoritesList.value.filter((item) =>{
                    return  item.id != id
                })
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
}

export const useUpdateFavorites = (item: GetFavoritesListItem ,createCollectDialogShow :Ref<boolean> ,updataInfo :Ref<GetFavoritesListItem>) => {
    updataInfo.value = item
    createCollectDialogShow.value = true
    console.log("renew", item)
}

export const useInit = async (favoritesList: Ref<GetFavoritesListRes> ,isLoading :Ref<boolean>) => {
    try {
        //Get list of favorites 
        const response = await getFavoritesList();
        if (!response.data) return false
        favoritesList.value = response.data
        isLoading.value = true

    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to obtain",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}
