<template>
  <div class="loading">
    <div class="block" v-for="i in 8">
      <div class="box"></div>
      <div class="line">
        <div class="line1"></div>
        <div class="line2"></div>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
@keyframes skeleton-loading {
  0% {
    background-position: 100% 50%;
  }

  100% {
    background-position: 0 50%;
  }
}

.loading {
  padding: 1.5em 0;

  .block {
    opacity: 0.5;
    margin-bottom: 1.5em;
    display: flex;
    justify-content: space-between;
    border-radius: 5px;
    overflow: hidden;

    .box,
    .line1,
    .line2 {
      // background-image: linear-gradient(90deg, #f2f2f2ee 25%, #e6e6e6 37%, #f2f2f2 63%);
      background-image: linear-gradient(
        90deg,
        var(--color-loading-b1) 25%,
        var(--color-loading-b2) 37%,
        var(--color-loading-b3) 63%,
        var(--color-loading-b1) 100%
      );
      background-size: 400% 100%;
      background-position: 100% 50%;
      animation: skeleton-loading 1.4s ease infinite;
    }

    .box {
      width: 18%;
      height: 6em;
      margin-bottom: 0;
    }

    .line {
      width: 80%;
    }

    .line1 {
      margin-bottom: 0.6em;
      height: 1em;
    }

    .line2 {
      height: 4.4em;
    }
  }
}

@media screen and (max-width: @phone) {
  .loading {
    .block {
      margin-bottom: 1em;

      .box {
        width: 100%;
        border-radius: 10px;
        height: 8em;
        margin-bottom: 0;
      }

      .line {
        display: none;
      }
    }
  }
}
</style>
