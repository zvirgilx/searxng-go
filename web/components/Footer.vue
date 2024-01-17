<template>
  <div class="footer-page">
    <footer>
      <p>some word</p>
    </footer>
  </div>
</template>

<style lang="less" scoped>
.footer-page {
  display: flex;
  justify-content: center;
  align-items: center;
  clear: both;
  min-height: 4rem;
  padding: 1em 0;
  text-align: center;
  background-color: var(--color-footer-background);
  border-top: 1px solid var(--color-footer-border);
  overflow: hidden;

  p {
    font-size: 0.8em;
  }
}
</style>
