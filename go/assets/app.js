(function () {
  'use strict';

  const root = document.documentElement;
  const course = root.dataset.course || 'course';
  const progressKey = 'study-' + course + '-progress-v1';
  const themeKey = 'study-theme';
  const toast = document.getElementById('toast');
  let toastTimer;

  function readStorage(key, fallback) {
    try {
      const value = localStorage.getItem(key);
      return value === null ? fallback : JSON.parse(value);
    } catch (_) {
      return fallback;
    }
  }

  function writeStorage(key, value) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch (_) {
      // 저장소를 쓸 수 없어도 문서는 그대로 읽을 수 있다.
    }
  }

  function showToast(message) {
    if (!toast) return;
    window.clearTimeout(toastTimer);
    toast.textContent = message;
    toast.classList.add('show');
    toastTimer = window.setTimeout(() => toast.classList.remove('show'), 1800);
  }

  /* 진도: 코스 전체 id 목록을 한 키에 저장하고, 페이지는 자기 체크박스만 갱신한다. */
  function readProgress() {
    const value = readStorage(progressKey, []);
    return new Set(Array.isArray(value) ? value.filter((item) => typeof item === 'string') : []);
  }

  const boxes = Array.from(document.querySelectorAll('input[data-progress]'));
  if (!boxes.length) document.querySelector('.toc .toc-progress')?.setAttribute('hidden', '');
  const cards = Array.from(document.querySelectorAll('.lesson-card[data-lesson]'));

  function setProgressText(prefix, done, total) {
    const percent = total ? Math.round((done / total) * 100) : 0;
    const bar = document.getElementById(prefix + '-bar');
    const count = document.getElementById(prefix + '-count');
    const total_ = document.getElementById(prefix + '-total');
    const pct = document.getElementById(prefix + '-percent');
    if (bar) {
      bar.value = percent;
      bar.textContent = percent + '%';
    }
    if (count) count.textContent = String(done);
    if (total_) total_.textContent = String(total);
    if (pct) pct.textContent = percent + '%';
  }

  function updateProgress() {
    const stored = readProgress();
    if (boxes.length) {
      boxes.forEach((box) => {
        box.checked = stored.has(box.dataset.progress);
      });
      const done = boxes.filter((box) => box.checked).length;
      setProgressText('lesson-progress', done, boxes.length);
    }
    if (cards.length) {
      let doneAll = 0;
      let totalAll = 0;
      cards.forEach((card) => {
        const total = Number(card.dataset.total) || 0;
        const prefix = course + '-' + card.dataset.lesson + '-';
        let done = 0;
        stored.forEach((id) => {
          if (id.startsWith(prefix)) done += 1;
        });
        done = Math.min(done, total);
        doneAll += done;
        totalAll += total;
        const bar = card.querySelector('progress');
        const label = card.querySelector('.count');
        if (bar) bar.value = total ? Math.round((done / total) * 100) : 0;
        if (label) label.textContent = done + '/' + total;
        card.classList.toggle('done', total > 0 && done >= total);
      });
      setProgressText('course-progress', doneAll, totalAll);
    }
  }

  boxes.forEach((box) => {
    box.addEventListener('change', () => {
      const stored = readProgress();
      if (box.checked) stored.add(box.dataset.progress);
      else stored.delete(box.dataset.progress);
      writeStorage(progressKey, Array.from(stored));
      updateProgress();
    });
  });

  document.getElementById('reset-progress')?.addEventListener('click', () => {
    if (!window.confirm('이 코스의 저장된 학습 진도를 모두 초기화할까요?')) return;
    writeStorage(progressKey, []);
    updateProgress();
    showToast('진도를 초기화했습니다.');
  });

  updateProgress();

  /* 테마 */
  const themeQuery = window.matchMedia?.('(prefers-color-scheme: dark)');
  const themeButton = document.getElementById('theme-toggle');
  const storedTheme = readStorage(themeKey, null);
  let themePreference = storedTheme === 'dark' || storedTheme === 'light' ? storedTheme : null;

  function preferredTheme() {
    return themePreference || (themeQuery?.matches ? 'dark' : 'light');
  }

  function applyTheme(theme) {
    const isDark = theme === 'dark';
    const label = isDark ? '밝은 테마로 전환' : '어두운 테마로 전환';
    root.dataset.theme = isDark ? 'dark' : 'light';
    themeButton?.setAttribute('aria-label', label);
    themeButton?.setAttribute('aria-pressed', String(isDark));
    themeButton?.setAttribute('title', label);
  }

  applyTheme(preferredTheme());
  themeButton?.addEventListener('click', () => {
    themePreference = root.dataset.theme === 'dark' ? 'light' : 'dark';
    writeStorage(themeKey, themePreference);
    applyTheme(themePreference);
  });
  themeQuery?.addEventListener?.('change', () => {
    if (!themePreference) applyTheme(preferredTheme());
  });

  /* 모바일 메뉴 */
  const topbar = document.querySelector('.topbar');
  const menuButton = document.querySelector('.menu-button');
  const nav = document.getElementById('primary-nav');

  function closeMenu(restoreFocus) {
    const wasOpen = nav?.classList.contains('open');
    nav?.classList.remove('open');
    menuButton?.setAttribute('aria-expanded', 'false');
    const label = menuButton?.querySelector('.sr-only');
    if (label) label.textContent = '메뉴 열기';
    if (restoreFocus && wasOpen) menuButton?.focus();
  }

  menuButton?.addEventListener('click', () => {
    const open = Boolean(nav?.classList.toggle('open'));
    menuButton.setAttribute('aria-expanded', String(open));
    const label = menuButton.querySelector('.sr-only');
    if (label) label.textContent = open ? '메뉴 닫기' : '메뉴 열기';
  });
  nav?.querySelectorAll('a').forEach((link) => link.addEventListener('click', () => closeMenu(false)));
  document.addEventListener('click', (event) => {
    if (nav?.classList.contains('open') && topbar && !topbar.contains(event.target)) closeMenu(false);
  });
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') closeMenu(true);
  });

  /* 코드 블록: 상단 막대(파일 경로·언어)와 복사 버튼을 자동으로 붙인다. */
  function legacyCopy(text) {
    if (typeof document.execCommand !== 'function') return false;
    const previousFocus = document.activeElement;
    const source = document.createElement('textarea');
    source.value = text;
    source.setAttribute('readonly', '');
    source.setAttribute('aria-hidden', 'true');
    source.className = 'copy-source';
    document.body.append(source);
    source.select();
    const copied = document.execCommand('copy');
    source.remove();
    previousFocus?.focus?.({ preventScroll: true });
    return copied;
  }

  async function copyText(text) {
    if (navigator.clipboard?.writeText) {
      try {
        await navigator.clipboard.writeText(text);
        return true;
      } catch (_) {
        // file:// 환경에서는 최신 클립보드 API가 막힐 수 있어 아래 방식으로 다시 시도한다.
      }
    }
    return legacyCopy(text);
  }

  document.querySelectorAll('pre.code').forEach((pre) => {
    const bar = document.createElement('div');
    const label = document.createElement('span');
    const button = document.createElement('button');
    const file = pre.dataset.src;
    bar.className = 'code-bar';
    label.className = 'label' + (file ? ' src' : '');
    label.textContent = file || pre.dataset.lang || 'code';
    button.type = 'button';
    button.className = 'copy-button';
    button.textContent = '복사';
    button.setAttribute('aria-label', (file || pre.dataset.lang || '코드') + ' 복사');
    button.addEventListener('click', async () => {
      const copied = await copyText(pre.textContent || '');
      showToast(copied ? '코드를 복사했습니다.' : '복사하지 못했습니다. 직접 선택해 주세요.');
    });
    bar.append(label, button);
    pre.classList.add('has-bar');
    pre.before(bar);
  });

  /* 레슨 목차와 현재 위치 표시 */
  const tocList = document.getElementById('toc-list');
  if (tocList) {
    const sections = Array.from(document.querySelectorAll('.lesson-section[id]'));
    sections.forEach((section) => {
      const heading = section.querySelector(':scope > h2');
      if (!heading) return;
      const clone = heading.cloneNode(true);
      clone.querySelectorAll('.tag').forEach((tag) => tag.remove());
      const item = document.createElement('li');
      const link = document.createElement('a');
      link.href = '#' + section.id;
      link.textContent = clone.textContent.trim();
      item.append(link);
      tocList.append(item);
    });

    if ('IntersectionObserver' in window) {
      const links = new Map(Array.from(tocList.querySelectorAll('a')).map((a) => [a.getAttribute('href').slice(1), a]));
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (!entry.isIntersecting) return;
            links.forEach((a) => a.classList.remove('active'));
            links.get(entry.target.id)?.classList.add('active');
          });
        },
        { rootMargin: '-80px 0px -65% 0px' },
      );
      sections.forEach((section) => observer.observe(section));
    }
  }

  /* 코스 홈 검색 */
  const searchInput = document.getElementById('course-search-input');
  const searchClear = document.getElementById('course-search-clear');
  const searchStatus = document.getElementById('course-search-status');
  const searchResults = document.getElementById('course-search-results');

  function collapse(value) {
    return value.replace(/\s+/g, ' ').trim();
  }

  function normalize(value) {
    return collapse(value.normalize('NFKC').toLocaleLowerCase('ko'));
  }

  const entries = cards.map((card) => {
    const link = card.querySelector('h3 a');
    return {
      href: link?.getAttribute('href') || '#',
      label: collapse(card.querySelector('.num')?.textContent || ''),
      title: collapse(link?.textContent || ''),
      summary: collapse(card.querySelector('p')?.textContent || ''),
      text: normalize((card.textContent || '') + ' ' + (card.dataset.keywords || '')),
    };
  });

  function resetSearch() {
    if (searchInput) searchInput.value = '';
    if (searchClear) searchClear.hidden = true;
    if (searchStatus) searchStatus.textContent = '검색어를 입력하세요.';
    if (searchResults) {
      searchResults.replaceChildren();
      searchResults.hidden = true;
    }
  }

  function renderSearch() {
    if (!searchInput || !searchStatus || !searchResults) return;
    const query = normalize(searchInput.value);
    if (searchClear) searchClear.hidden = !query;
    searchResults.replaceChildren();
    if (!query) {
      resetSearch();
      return;
    }
    const tokens = query.split(' ');
    const matches = entries.filter((entry) => tokens.every((token) => entry.text.includes(token)));
    searchResults.hidden = false;
    if (!matches.length) {
      const empty = document.createElement('li');
      empty.className = 'search-empty';
      empty.textContent = '일치하는 레슨이 없습니다. 더 짧거나 다른 용어로 검색해 보세요.';
      searchResults.append(empty);
      searchStatus.textContent = '검색 결과 0개';
      return;
    }
    matches.slice(0, 8).forEach((entry) => {
      const item = document.createElement('li');
      const link = document.createElement('a');
      const label = document.createElement('span');
      const title = document.createElement('strong');
      const summary = document.createElement('small');
      link.href = entry.href;
      label.textContent = entry.label;
      title.textContent = entry.title;
      summary.textContent = entry.summary;
      link.append(label, title, summary);
      item.append(link);
      searchResults.append(item);
    });
    searchStatus.textContent = matches.length > 8 ? '총 ' + matches.length + '개 중 8개 표시' : '검색 결과 ' + matches.length + '개';
  }

  searchInput?.addEventListener('input', renderSearch);
  searchInput?.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && searchInput.value) {
      event.stopPropagation();
      resetSearch();
    }
  });
  searchClear?.addEventListener('click', () => {
    resetSearch();
    searchInput?.focus();
  });

  window.addEventListener('storage', (event) => {
    if (event.key === progressKey || event.key === null) updateProgress();
    if (event.key === themeKey || event.key === null) {
      const next = readStorage(themeKey, null);
      themePreference = next === 'dark' || next === 'light' ? next : null;
      applyTheme(preferredTheme());
    }
  });
})();
