<template>
  <div class="wrapper index-page">
    <main id="main_index">
      <div class="title">
        <h1>SearXNG Internal</h1>
      </div>

      <SearchBox :isLogoShow="false"></SearchBox>
    </main>

    <Footer></Footer>
  </div>
</template>

<style lang="less" scoped>
.title {
  text-align: center;
  font-size: 2em;
  margin-bottom: 1em;
}

.index-page #main_index {
  margin-top: 20vh;
}

@media screen and (max-width: @phone) {

  body.search-input-focus {
    .title {
      display: none;
    }

    .index-page #main_index {
      margin-top: 10px;
    }
  }

}
</style>
