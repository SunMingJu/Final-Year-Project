<template>
  <div class="side-navigation">
    <div>
      <div class="avatar-box">
        <el-avatar :size="110" :src="userInfo.userInfoData.photo" />
      </div>
      <p>{{ userInfo.userInfoData.username }}</p>
    </div>
    <el-menu class="el-menu-vertical" :collapse-transition="false" @open="handleOpen" @close="handleClose"
      @select="handleSelect">
      <el-sub-menu index="UserShow">
        <template #title>
          <SvgIcon name="user" class="icon"></SvgIcon>
          <span>User Info</span>
        </template>
        <el-menu-item index="UserInfo">
          <div class="icon-item">
            <SvgIcon name="userData" class="icon-small"></SvgIcon>
            <span>personal information</span>
          </div>
        </el-menu-item>
        <el-menu-item index="PictureSetting">
          <div class="icon-item">
            <SvgIcon name="pictures" class="icon-small"></SvgIcon>
            <span>Avatar settings</span>
          </div>
        </el-menu-item>
        <el-menu-item index="Safety">
          <div class="icon-item">
            <SvgIcon name="security" class="icon-small"></SvgIcon>
            <span>User security</span>
          </div>
        </el-menu-item>
      </el-sub-menu>
      <el-menu-item index="LiveSetUp">
        <div class="icon-item">
          <SvgIcon name="live" class="icon"></SvgIcon>
          <span>Live broadcast settings</span>
        </div>
      </el-menu-item>
      <el-sub-menu index="Favorites">
        <template #title>
          <SvgIcon name="collection" class="icon"></SvgIcon>
          <span>my collection</span>
        </template>
        <el-menu-item @click="createCollectDialogShow = true" index="CreateFavorites">
          <div class="icon-item">
            <SvgIcon name="folder" class="icon-smallxl"></SvgIcon>
            <span>Create favorites</span>
          </div>
        </el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="Record">
        <template #title>
          <SvgIcon name="playRecording" class="icon"></SvgIcon>
          <span>history record</span>
        </template>
        <el-menu-item @click="clearHistory()" index="ClearHistory">
          <div class="icon-item">
            <SvgIcon name="trashCan" class="icon"></SvgIcon>
            <span>Clear history</span>
          </div>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>

    <!--Message pop-up box -->
    <div class="dialog">
      <el-dialog v-model="createCollectDialogShow" title="Create favorites" width="50%" center :close-on-click-modal="true"
        :append-to-body="true" :before-close="createCollectDialogShowClose" align-center>
        <CreateFavorite :type="true" @shutDown="shutDown"></CreateFavorite>
      </el-dialog>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/store/main";
import { useRouter } from "vue-router"
import CreateFavorite from "@/components/createFavorites/createFavorites.vue"
import { ref } from "vue";
import Swal from "sweetalert2";
import globalScss from "@/assets/styles/global/export.module.scss"
import { clearRecord } from "@/apis/personal";

comments: {
  CreateFavorite
}
const userInfo = useUserStore();
const router = useRouter();

const createCollectDialogShow = ref(false)

const handleOpen = (key: string, keyPath: string[]) => {
  router.push({ name: key })
};

const handleClose = (key: string, keyPath: string[]) => {
  router.push({ name: key })
};

const handleSelect = (key: string, keyPath: string[]) => {
  let arr = ["CreateFavorites", "ClearHistory"]
  if (key == "" || arr.indexOf(key) >= 0) return false
  router.push({ name: key })
}

const createCollectDialogShowClose = (done: () => void) => {
  done()
}

const clearHistory = () => {
  Swal.fire({
    heightAuto: false,
    title: 'Confirm to clear history?',
    icon: 'warning',
    confirmButtonColor: globalScss.colorButtonTheme,
    showCancelButton: true,
    confirmButtonText: 'confirm',
    cancelButtonText: 'Cancel',
    reverseButtons: true
  }).then(async (result) => {
    if (result.isConfirmed) {
      try {
        await clearRecord()
        Swal.fire({
          title: "Cleared successfully",
          confirmButtonColor: globalScss.colorButtonTheme,
          heightAuto: false,
          icon: "success",
          preConfirm: () => {
            router.push({ name: "Record", query: { type: 'createTime' + Date.now() } })
          }
        })
      } catch (err: any) {
        console.log(err)
        Swal.fire({
          title: "Clearing failed",
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
//Close create collect dialog show
const shutDown = () => {
  createCollectDialogShow.value = false
}

</script>

<style scoped lang="scss">
@import "./static/style/personalNavigation.scss";
</style>
