<template>
    <div class="carousel">
        <el-carousel ref="carouselRef" arrow="never" height="490px" :autoplay="true" @change="change">
            <el-carousel-item v-for="item in rotograph" :key="item.cover">
                <div class="carousel-item">
                    <div class="carousel-img-box">
                        <el-image class="carousel-img" :src="item.cover" fit="cover" />
                    </div>
                </div>
            </el-carousel-item>
            <div class="gradient" :style="{ background: `linear-gradient(rgba(0,0,0,0),${color})` }"></div>
            <!--Bottom of the carousel -->
            <div class="carousel-bottom" :style="{ backgroundColor: `${color}` }">
                <div class="carousel-title">{{ carouselTitle ? carouselTitle : props.rotograph[0]?.title }}</div>
                <div class="toggle-button">
                    <div class="button-item" v-throttle="{ fun: carouselSwitch, params: [false], time: 500 }">
                        <SvgIcon name="rightArrow" class="arrow-icon rotation" color="#fff">
                        </SvgIcon>
                    </div>
                    <div class="button-item" v-throttle="{ fun: carouselSwitch, params: [true], time: 500 }">
                        <SvgIcon name="leftArrow" class="arrow-icon" color="#fff"></SvgIcon>
                    </div>
                </div>
            </div>
        </el-carousel>
    </div>
</template>

<script lang="ts" setup>
import { RotographList } from "@/types/home/home";
import { vThrottle } from "@/utils/customInstruction/throttle";
import { ElCarousel, ElCarouselItem, ElImage } from 'element-plus';
import { onMounted, ref } from 'vue';

const props = defineProps({
    rotograph: {
        type: Array as () => RotographList,
        required: true,
        default: () => []
    }
})

const carouselRef = ref()
const carouselTitle = ref("")
const color = ref("")
//Switch carousel image true next false previous
const carouselSwitch = (is: boolean) => {
    if (is) {
        carouselRef.value.next()

    } else {
        carouselRef.value.prev()
    }
}


//Carousel chart change event
const change = (index: number) => {
    carouselTitle.value = props.rotograph[index].title
    color.value = props.rotograph[index].color
}

onMounted(() => {
    color.value = props.rotograph[0].color
})
</script>

<style lang="scss" scoped>
@import "./static/css/homeRotograph.scss";
</style>