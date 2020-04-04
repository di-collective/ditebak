<style>
  small {
    display: block;
    line-height: 1.5rem;
  }

  .footer {
    background-color: var(--mdc-theme-primary,#fc270f);
    color: white;
    height: 10rem;
  }
  .grid-container {
    display: grid;
    grid-template-columns: auto auto;
  }
  .text-r {
    text-align: right;
  }
  .social {
    color: #FC270F;
    background: white;
    padding: 0.1rem 0.5rem;
    border-radius: 10px;
    -webkit-box-shadow: 2px 2px 5px 0px rgba(87,87,87,1);
    -moz-box-shadow: 2px 2px 5px 0px rgba(87,87,87,1);
    box-shadow: 2px 2px 5px 0px rgba(87,87,87,1);
  }
</style>

<footer class="footer">
  <div class="grid-container pad-t pad-r pad-b pad-l mobile-wd">
    <div>
      <a href="https://ko-fi.com/O5O81KSO7" target="_blank">
        <img height="124" border="0" src="/ko-fi.png" alt="Buy Me a Coffee at ko-fi.com">
      </a>
    </div>
    <div class="text-r">
      <small>Copyright &copy; ditebak.com {new Date().getFullYear()}</small>
      <small><a class="social" target="_blank" href="https://github.com/bastianrob">personal github</a></small>
      <!-- <small><a class="social" target="_blank" href="https://github.com/di-collective/ditebak">ditebak github</a></small> -->
      <small><a class="social" target="_blank" href="https://medium.com/@BastianRob">medium blog</a></small>
      <small><a class="social" target="_blank" href="https://dev.to/bastianrob">dev.to blog</a></small>
    </div>
  </div>
</footer>
