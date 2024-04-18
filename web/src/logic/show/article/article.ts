
import { Router, useRouter, useRoute, RouteLocationNormalizedLoaded } from 'vue-router';
import { getArticleComment, getArticleContributionByID } from "@/apis/contribution"
import { GetArticleCommentReq, GetArticleContributionByIDReq, GetArticleContributionByIDRes } from "@/types/show/article/article"
import { Ref, reactive, ref } from "vue"
import globalScss from "@/assets/styles/global/export.module.scss"
import Swal from 'sweetalert2';
import { Size } from 'tsparticles-engine';

export const useArticleShowProp = () => {
    const articleID = ref(0)
    const articleInfo = ref(<GetArticleContributionByIDRes>{})
    const router = useRouter()
    const route = useRoute()

    //Reply to secondary comments
    const replyCommentsDialog = reactive({
        show: false,
        commentsID: 0,
    })

    return {
        articleID,
        articleInfo,
        router,
        route,
        replyCommentsDialog,
    }
}
export const useInit = async (articleID: Ref<number>, articleInfo: Ref<GetArticleContributionByIDRes>, route: RouteLocationNormalizedLoaded, router: Router) => {
    try {
        if (!route.query.articleID) {
            router.back()
            Swal.fire({
                title: "Missing article id",
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
                icon: "error",
            })
            router.back()
            return
        }
        articleID.value = Number(route.query.articleID)
        window.onresize = function () {
            const canvasSnow = document.getElementById('canvas_sakura') as HTMLEmbedElement;
            if (!canvasSnow) return false
            canvasSnow.width = String(window.innerWidth);
            canvasSnow.height = String(window.innerHeight);
        }
        //Get article content
        const response = await getArticleContributionByID(<GetArticleContributionByIDReq>{
            articleID: articleID.value
        })

        if (!response.data) throw "The article content is empty"
        articleInfo.value = response.data
        articleInfo.value.comments = response.data?.comments

    } catch (err) {
        console.log(err)
    }
}
