<template>
  <div class="wrapper search-page">
    <main id="main_index">
      <SearchBox></SearchBox>

      <div id="results">
        <div class="urls">
          <div v-if="loading">
            <SearchLoading></SearchLoading>
          </div>
          <div v-else>
            <template v-if="results.length > 0">
              <article class="result" v-for="item in results" :key="item.url">
                <a
                  target="_blank"
                  v-if="item.img_src || item.thumbnail"
                  :href="item.url"
                  rel="noreferrer"
                >
                  <img class="image" :src="item.img_src || item.thumbnail" :title="item.title" />
                </a>
                <h3>
                  <a target="_blank" :href="item.url" rel="noreferrer">
                    <span v-html="getHighlightedText(item.title, q)"></span>
                  </a>
                </h3>
                <p class="content">
                  {{ item.content }}
                </p>
                <div class="engines">
                  <span>{{ item.engine }}</span>
                  <a
                    target="_blank"
                    :href="getCacheUrl(item.url)"
                    class="cache_link"
                    rel="noreferrer"
                  >
                    cached
                  </a>
                </div>
              </article>
            </template>
            <template v-else>
              <div class="results-empty">
                <p><strong>Sorry!</strong></p>
                <p>No results were found. You can try to:</p>
                <ul>
                  <li>Refresh the page.</li>
                  <li>Search for another query or select another category (above).</li>
                </ul>
              </div>
            </template>
          </div>
        </div>

        <!-- <div id="backToTop">
          <a href="javascript:;"><svg class="ion-icon-small" viewBox="0 0 512 512" aria-hidden="true">
              <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="48"
                d="M112 328l144-144 144 144"></path>
            </svg></a>
        </div> -->

        <div v-if="results.length !== 0">
          <Pagination :total="total"></Pagination>
        </div>
      </div>
    </main>

    <Footer></Footer>
  </div>
</template>

<script setup>
import { getHighlightedText } from "~/utils";
const { routeChangeListener, getPageNo } = useCustomRouter();

const results = ref([]);
const total = ref(0);
const loading = ref(false);
const q = ref("");

routeChangeListener(() => {
  searchInit();
});

function getCacheUrl(url) {
  return `https://web.archive.org/web/${url}`;
}

async function searchInit() {
  console.log("call searchInit");

  const route = useRoute();

  loading.value = true;
  results.value = [];

  q.value = route.query.q ?? "";

  const page_no = getPageNo();
  const data = await useFetchSearch(q, page_no);

  try {
    console.log("data", data, data.value, data.value.results, typeof data.value.results);

    if (data?.value?.results) {
      results.value = data.value.results;

      // no total, so set 100
      total.value = 100;
    }
  } catch (err) {
    console.error(err);
  }

  loading.value = false;
}
</script>

<style lang="less">
.search-page {
  margin-top: 1.5em;
}

#results {
  width: 46rem;
  margin: 0 auto;
}

.result {
  margin: @results-margin 0;
  padding: @result-padding;
  // .ltr-border-left(0.2rem solid transparent);

  h3 {
    font-size: 1.1rem;
    line-height: 1.2;
    word-wrap: break-word;
    margin: 0.4rem 0 0.4rem 0;
    padding: 0;

    a {
      color: var(--color-result-link-font);
      font-weight: normal;

      &:visited {
        color: var(--color-result-link-visited-font);
      }

      &:focus,
      &:hover {
        text-decoration: underline;
        border: none;
        outline: none;
      }

      span {
        em {
          font-weight: bold;
          color: #f73131;
        }
      }
    }
  }

  .cache_link,
  .proxyfied_link {
    font-size: 0.9em !important;
  }

  .content {
    font-size: 0.9em;
    margin: 0;
    padding: 0;
    max-width: 54em;
    word-wrap: break-word;
    line-height: 1.24;

    // overflow: hidden;
    // text-overflow: ellipsis;
    // display: -webkit-box;
    // -webkit-line-clamp: 5;
    // -webkit-box-orient: vertical;

    .highlight {
      color: var(--color-result-description-highlight-font);
      background: inherit;
      font-weight: bold;
    }
  }

  img {
    &.image {
      float: left;
      margin: 0.5rem 1rem 0 0;
      width: 7rem;
      max-height: 7rem;
      object-fit: scale-down;
      object-position: center;
    }
  }

  .engines {
    clear: both;
    text-align: right;
    line-height: 1;
    margin-top: 1.5em;
    font-size: 0.8em;

    span {
      margin-right: 0.5em;
    }
  }
}

#backToTop {
  border: 1px solid var(--color-backtotop-border);
  margin: 0;
  padding: 0;
  font-size: 1em;
  background: var(--color-backtotop-background);
  position: fixed;
  bottom: 8rem;
  .ltr-left(@results-width + @results-offset + (0.5 * @results-gap - 1.2em));
  transition: opacity 0.5s;
  opacity: 1;
  pointer-events: none;
  .rounded-corners;

  a {
    display: block;
    margin: 0;
    padding: 0.7em;
  }

  a,
  a:visited,
  a:hover,
  a:active {
    color: var(--color-backtotop-font);
  }
}

.results-empty {
  position: relative;
  display: flex;
  padding: 1rem;
  margin: 2em 0 1em 0;
  border: 1px solid var(--color-toolkit-dialog-border);
  text-align: left;
  border-radius: 10px;
  display: block;
  color: var(--color-error);
  background: var(--color-error-background);
  border-color: var(--color-error);
  line-height: 1.4;

  ul {
    padding-left: 2em;

    li {
      list-style: disc;
    }
  }
}

/*
  phone layout
*/

@media screen and (max-width: @phone) {
  html {
    background-color: var(--color-base-background-mobile);
  }

  .search-page {
    margin-top: 0;
  }

  #backToTop {
    display: none;
  }

  #results {
    width: auto;
  }

  .result {
    background: var(--color-result-background);
    border: 1px solid var(--color-result-background);
    margin: 1rem 0;
    .rounded-corners;
  }
}

/*
  small-phone layout
*/

// @media screen and (max-width: @small-phone) {
//   .result-videos {
//     img.thumbnail {
//       float: none !important;
//     }

//     .content {
//       overflow: inherit;
//     }
//   }
// }

// pre code {
//   white-space: pre-wrap;
// }
</style>
