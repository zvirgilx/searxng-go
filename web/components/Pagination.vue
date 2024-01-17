<template>
  <div class="pagination">
    <div class="prev-page">
      <button @click="goPrevPage()" v-if="curPage !== 1" role="link" type="button">
        <svg class="ion-icon-small" viewBox="0 0 512 512" aria-hidden="true">
          <path
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="48"
            d="M328 112L184 256l144 144"
          ></path>
        </svg>
        Previous page
      </button>
    </div>

    <div class="number-page">
      <ol>
        <li v-for="page in pages">
          <span v-if="page === curPage">{{ page }}</span>
          <a v-else href="javascript:;" @click="goCurPage(page)">{{ page }}</a>
        </li>
      </ol>
    </div>

    <div class="next-page">
      <button @click="goNextPage()" v-if="curPage < total" role="link" type="button">
        Next page
        <svg class="ion-icon-small" viewBox="0 0 512 512" aria-hidden="true">
          <path
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="48"
            d="M184 112l144 144-144 144"
          ></path>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { useCustomRouter } from "~/utils/hooks";
const { goPage, routeChangeListener } = useCustomRouter();

const props = defineProps({
  total: {
    type: Number,
    default: 1,
  },
});

const pages = ref([]);
let curPage = 1;

routeChangeListener(() => {
  pageInit();
});

function pageInit() {
  const route = useRoute();
  console.log("route.query", route.query);

  // const q = route.query.q ?? ""; // search keyword
  curPage = route.query.page_no ? parseInt(route.query.page_no, 10) : 1; // page number

  pages.value = getPagesNumberArrByCurPage(curPage);

  console.log(curPage, pages.value);
}

function scrollTop() {
  window?.scrollTo({
    top: 0,
    behavior: "instant",
    // behavior: "smooth"
  });
}

function goCurPage(page) {
  scrollTop();
  goPage(page);
}

function goPrevPage() {
  scrollTop();
  goPage(curPage - 1);
}

function goNextPage() {
  scrollTop();
  goPage(curPage + 1);
}

function getPagesNumberArrByCurPage(curPage) {
  const ret = [];

  const step = 3;
  const middle = 1 + 1 * step;
  const large = 1 + 2 * step;

  if (props.total <= large) {
    for (let i = 1; i <= props.total; i++) {
      ret.push(i);
    }
  } else {
    if (curPage < middle) {
      for (let i = 1; i <= large; i++) {
        ret.push(i);
      }
    } else {
      for (let i = curPage - step; i <= curPage + step; i++) {
        if (i >= 1 && i <= props.total) {
          ret.push(i);
        }
      }
    }
  }

  return ret;
}
</script>

<style lang="less" scoped>
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 2rem 0 0;

  span,
  a,
  button {
    background: var(--color-pagination-background);
    color: var(--color-pagination-font);
  }

  button:hover,
  a:hover {
    background: var(--color-pagination-select-background);
    color: var(--color-pagination-select-font);
  }

  button {
    padding: 0.7rem;
    display: flex;

    border-radius: 10px;
    border: 0;
    cursor: pointer;
  }

  .number-page {
    ol {
      list-style: none;
      display: flex;

      li {
        margin: 0 0.5em;

        span,
        a {
          width: 36px;
          height: 36px;
          line-height: 36px;
          border: none;
          border-radius: 6px;
          display: inline-block;
          vertical-align: text-bottom;
          text-align: center;
          text-decoration: none;
        }

        span {
          color: var(--color-pagination-select-font);
          background: var(--color-pagination-select-background);
        }
      }
    }
  }
}

@media screen and (max-width: @phone) {
  .pagination {
    .number-page {
      display: none;
    }
  }
}
</style>
