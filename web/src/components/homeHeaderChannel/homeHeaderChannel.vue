<template>
       <!--Top channel popular categories -->
       <div class="header-channel">
      <!--Left part -->
      <div class="header-channel-left">
        <div class="channel-left-item">
          <div class="channel-left-icon bg-orange">
            <SvgIcon name="dynamic" class="icon" color="#fff"></SvgIcon>
          </div>
          <div class="channel-left-item-title">dynamic</div>
        </div>
        <div class="channel-left-item">
          <div class="channel-left-icon bg-pink">
            <SvgIcon name="flame" class="icon" color="#fff"></SvgIcon>
          </div>
          <div class="channel-left-item-title">Popular</div>
        </div>
        <div class="channel-left-item">
          <div class="channel-left-icon bg-green">
            <SvgIcon name="channel" class="icon" color="#fff"></SvgIcon>
          </div>
          <div class="channel-left-item-title">channel</div>
        </div>
      </div>
      <!--middle part -->
      <div class="header-channel-right">
        <div class="channel-right-box-left">
          <div class="right-box-left-item" v-for="item in 24">
            Fan drama
          </div>
        </div>

        <div class="channel-right-box-right">
          <div class="right-box-right-item" v-for="sideMavigationItem in sideMavigation">
            <SvgIcon :name="sideMavigationItem.iconName" class="icon" color="#61666D"></SvgIcon>
            <span>{{ sideMavigationItem.title }}</span>
          </div>
        </div>
      </div>

    </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';

const sideMavigation = reactive([
    {iconName : "columnInformation" , title :"Column"},
    {iconName : "liveStreamPlayer" , title :"live streaming"},
    {iconName : "flag" , title :"Activity"},
    {iconName : "classroom" , title :"classroom"},
    {iconName : "chatList" , title :"Message list"},
    {iconName : "music" , title :"listen to a song"},
])

 
</script>

<style scoped lang="scss">
@import "./static/css/homeHeaderChannel.scss";
</style>