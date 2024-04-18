<template>
    <div class="comments-info">
        <textarea class="comments-textarea" placeholder="写下点什么..." maxlength="1000"
            v-model="comments.comments"></textarea>
    </div>
    <!--Labeling and publishing -->
    <div class="comments-fun">
        <div>
            <SvgIcon name="expression"
                :class="{ 'icon-small': true, 'animate__animated': true, 'animate__tada': emoji.animation }"
                @click="clickEmoji" @mouseleave="emoji.animation = false" @mouseover="emoji.animation = true">
            </SvgIcon>
        </div>

        <div> <el-button type="primary" v-removeFocus round @click="postComment(comments, articleID)">release</el-button>
        </div>
    </div>
    <div :class="{ 'comments-emoji': true, 'animate__animated': true, 'animate__headShake': !emoji.animationBox, 'animate__flipOutX': emoji.animationBox }"
        v-show="emoji.show">
        <span class="emoji-item" v-for="emojiItem in emoji.emoji" :key="emojiItem"
            @click="comments.comments = comments.comments + emojiItem">{{ emojiItem }}</span>
    </div>
</template>
<script lang="ts" setup>import { articlePostComment, getArticleComment } from '@/apis/contribution';
import { ArticlePostCommentReq, CommentsInfo, GetArticleCommentReq } from '@/types/show/article/article';
import { ElButton } from 'element-plus';
import Swal from 'sweetalert2';
import { vRemoveFocus } from "@/utils/customInstruction/focus"
import globalScss from "@/assets/styles/global/export.module.scss"
import { Ref, UnwrapNestedRefs, reactive, toRefs } from 'vue';

const props = defineProps({
    articleID: {
        type: Number,
        required: true,
    },
    articleInfo: {
        type: Object,
        required: true,
    },
    commentsID: {
        type: Number,
        required: true,
    }
})
const emit = defineEmits(['updateArticleInfo'])



const emoji = reactive({
    show: false,
    animation: false,
    animationBox: true,
    emoji: [
        "😀", "😄", "😁", "😆", "😅", "🤣", "😂", "🙂", "🙃", "😉", "😊", "😇", "🥰", "😍", "🤩", "😚", "🤗", "🤨", "😐", "😑", "😶", "🤐", "😏", "😒", "😮‍💨", "🤥", "😌", "😪", "🤤", "😷", "🤒", "🤕", "🥵", "😵",
        "😕", "🙁", "☹️", "😳", "😞", "😭", "🥱", "😩", "😰", "😲", "😯", "😠", "😩", "😧", "😯", "🥺"
    ]
})

const comments = reactive(<CommentsInfo>{
    comments: "",
})


//Click on emoticon to trigger
const clickEmoji = () => {
    emoji.animationBox = !emoji.animationBox
    if (emoji.show) {
        //Retract panel animation effect
        setTimeout(() => {
            emoji.show = !emoji.show
        }, 800);
    } else {
        emoji.show = !emoji.show
    }
}

//Reply to comment
const postComment = async (comments: UnwrapNestedRefs<CommentsInfo>, articleID: number) => {
    try {
        const requistData = <ArticlePostCommentReq>{
            article_id: articleID,
            content: comments.comments,
            content_id: props.commentsID,
        }
        const reponse = await articlePostComment(requistData)
        console.log(reponse)

        //Clear input money
        comments.comments = ""
        const commentsList = await getArticleComment(<GetArticleCommentReq>{ articleID: articleID })
        if (!commentsList.data) {
            throw ("Get Comment failed")
        }
        emit('updateArticleInfo', commentsList.data)
        const Toast = Swal.mixin({
            toast: true,
            position: 'top',
            showConfirmButton: false,
            timer: 3000,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({
            icon: 'success',
            title: 'Comment successful'
        })

    } catch (err) {
        Swal.fire({
            title: err as string,
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}

</script>

<style scoped lang="scss">
@import "./static/style/commentPosting.scss"
</style>
