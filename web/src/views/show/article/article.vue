<template>
  <div class="article-page">
    <div class="head">
      <topNavigation color="#fff" scroll :displaySearch="false"></topNavigation>
    </div>
    <!--Cover image -->
    <div class="cover-picture" :style="{ backgroundImage: `url(${articleInfo.cover})` }">
      <div class="article-info-container">
        <div class="title">{{ articleInfo.title }}</div>
        <div class="article-info">
          <span class="info-lb">
            <SvgIcon name="camera" class="icon-small"></SvgIcon> {{
              dayjs(articleInfo.created_at).format('YYYY.MM.DD.hh.mm')
            }}
          </span>
          <span class="info-lb">
            <SvgIcon name="hot" class="icon-small"></SvgIcon>
            <span>
              {{ articleInfo.heat }}
            </span>
          </span>
          <span class="info-lb">
            <SvgIcon name="comments" class="icon-small"></SvgIcon>
            <span>
              {{ articleInfo.comments_number }}
            </span>
          </span>
          <span class="info-lb">
            <SvgIcon name="like" class="icon-small"></SvgIcon>
            <span>
              {{ articleInfo.likes_number }}
            </span>
          </span>
        </div>
      </div>
    </div>
    <!--Article content -->
    <div class="article-content">
      <div class="content">
        <!--Other ql-container ql-snow -->
        <div class="ql-editor" v-html="articleInfo.content">
        </div>
      </div>
      <!--Bottom comments -->
      <div class="comments-box">
        <div class="comments-head">
          <SvgIcon name="editor" class="icon-edit"></SvgIcon> Comment
        </div>
        <div class="comments-main">
          <commentPosting :articleID="articleID" :articleInfo="articleInfo" @updateArticleInfo="updateArticleInfo"
            :commentsID="0"></commentPosting>
        </div>
        <div class="comments-show">
          <div class="comments-show-titel"><span>Comments | </span> <span>{{ articleInfo.comments_number }} comments</span>
          </div>
          <div class="comments-show-info">
            <div class="comment-info-detail" v-for="commentsItem in articleInfo.comments" :key="commentsItem.id">
              <el-avatar shape="square" :size="36" :src="commentsItem.photo" />
              <div class="comment-info-content">
                <div class="content-head">
                  <div> <span class="comment-info-username">{{ commentsItem.username }}</span> <span
                      class="comment-info-other">{{ dayjs(commentsItem.created_at).format('YYYY.MM.DD.hh.mm') }}</span>
                  </div>
                  <div class="commentInfo-reply">
                    <el-button type="primary" v-removeFocus round size="small"
                      @click="replyComments(commentsItem.id)">reply</el-button>
                  </div>
                </div>
                <!-- 评论内容部分 -->
                <div class="content-info">
                  {{ commentsItem.context }}
                </div>
                <!--Subcomment -->
                <div class="comment-info-detail" v-for="lowerComments in commentsItem.lowerComments"
                  :key="lowerComments.id">
                  <el-avatar shape="square" :size="36" :src="lowerComments.photo" />
                  <div class="comment-info-content">
                    <div class="content-head">
                      <div> <span class="comment-info-username">{{ lowerComments.username }}</span> <span
                          class="comment-info-other">{{ dayjs(lowerComments.created_at).format('YYYY.MM.DD.hh.mm')
                          }}</span>
                      </div>
                      <div class="commentInfo-reply">
                        <el-tag effect="dark" round v-removeFocus @click="replyComments(lowerComments.id)">
                          reply
                        </el-tag>
                      </div>
                    </div>
                    <!--Comment content part -->
                    <div class="content-info">
                      <span v-if="lowerComments.comment_user_id != commentsItem.uid"><span class="at-user">@{{
                        lowerComments.comment_user_name
                      }} </span> : </span> {{ lowerComments.context }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!--Reply to comment dialog-->
        <el-dialog v-model="replyCommentsDialog.show" title="Shipping address">
          <commentPosting :articleID="articleID" :articleInfo="articleInfo" @updateArticleInfo="updateArticleInfo"
            :commentsID="replyCommentsDialog.commentsID"></commentPosting>
        </el-dialog>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts" >
import commentPosting from "@/components/commentPosting/commentPosting.vue"

import 'quill/dist/quill.bubble.css'
import 'quill/dist/quill.core.css'
import 'quill/dist/quill.snow.css'


import topNavigation from "@/components/topNavigation/topNavigation.vue"
import { useArticleShowProp, useInit } from "@/logic/show/article/article"
import dayjs from "dayjs"
//Code highlighting
import { GetArticleCommentRes } from "@/types/show/article/article"
import { vRemoveFocus } from "@/utils/customInstruction/focus"
import { blossom } from "@/utils/effect/blossom"
import 'highlight.js/styles/agate.css'
import { onMounted, onUnmounted } from "vue"

components: {
  topNavigation
  commentPosting
}

const { articleID, articleInfo, router, route, replyCommentsDialog } = useArticleShowProp()
//Update comment data
const updateArticleInfo = (commentsList: GetArticleCommentRes) => {
  articleInfo.value.comments = commentsList.comments
  articleInfo.value.comments_number = commentsList.comments_number
}

//Reply to secondary comments
const replyComments = (commentsID: number) => {
  replyCommentsDialog.commentsID = commentsID
  replyCommentsDialog.show = !replyCommentsDialog.show
}

//Sakura animation
const { startSakura, stopp } = blossom()

onMounted(async () => {
  startSakura()
  await useInit(articleID, articleInfo, route, router)
})

onUnmounted(() => {
  stopp()
})

</script>

<style scoped lang="scss">
//Add related styles
@import "@/assets/styles/views/show/article/article.scss";
</style>
