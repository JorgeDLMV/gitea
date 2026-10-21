import {execPseudoSelectorCommands, handleFetchActionErrorFields, handleFetchActionSuccessJson, performFetchAction} from './fetch-action.ts';
import {GET, request} from './fetch.ts';
import {createElementFromHTML} from '../utils/dom.ts';
import {captureNavigations, normalizeTestHtml} from '../utils/testhelper.ts';

vi.mock('./fetch.ts', () => ({GET: vi.fn(), request: vi.fn()}));

test('execPseudoSelectorCommands', () => {
  window.document.body.innerHTML = `
<div id="d1">
    <ul id="u1">
        <li class="x"></li>
    </ul>
    <ul id="u2">
        <li class="x"></li>
    </ul>
</div>
<div id="d2">
    <ul id="u3">
        <li class="x"></li>
    </ul>
</div>`;

  let ret = execPseudoSelectorCommands(document.querySelector('#u1')!, '');
  expect(ret.targets).toEqual([document.querySelector('#u1')]);

  ret = execPseudoSelectorCommands(document.querySelector('#u1')!, '$this');
  expect(ret.targets).toEqual([document.querySelector('#u1')]);
  expect(ret.cmdInnerHTML).toBeFalsy();
  expect(ret.cmdMorph).toBeFalsy();

  ret = execPseudoSelectorCommands(document.querySelector('#u1')!, '$body $morph $innerHTML');
  expect(ret.targets).toEqual([document.body]);
  expect(ret.cmdInnerHTML).toBeTruthy();
  expect(ret.cmdMorph).toBeTruthy();

  ret = execPseudoSelectorCommands(document.querySelector('#u1')!, '$body .x');
  expect(ret.targets.length).toEqual(3);
  expect(ret.targets).toEqual(Array.from(document.querySelectorAll('.x')));

  ret = execPseudoSelectorCommands(document.querySelector('#u1 .x')!, '$closest(div) .x');
  expect(ret.targets.length).toEqual(2);
  expect(ret.targets).toEqual(Array.from(document.querySelectorAll('#d1 .x')));
});

test('handleFetchActionSuccessJson', () => {
  const navigations = captureNavigations();
  handleFetchActionSuccessJson(document.body, {redirect: '/'});
  handleFetchActionSuccessJson(document.body, {redirect: ''});
  handleFetchActionSuccessJson(document.body, {});
  expect(navigations.map((n) => n.type)).toEqual(['push', 'reload', 'reload']);
});

test('performFetchAction successReload', async () => {
  window.document.body.innerHTML = `<div id="a">old a</div><div id="b">old b</div><div id="c">old c</div>`;
  vi.mocked(request).mockResolvedValue(new Response('{}', {headers: {'Content-Type': 'application/json'}}));
  vi.mocked(GET).mockResolvedValue(new Response(`<div id="a">new a</div><div id="b">new b</div><div id="c">new c</div>`));
  const navigations = captureNavigations();
  await performFetchAction(document.body, {method: 'POST', url: '/action', successReload: '#a, #b'});
  expect(navigations).toEqual([]);
  expect(document.body.innerHTML).toEqual(`<div id="a">new a</div><div id="b">new b</div><div id="c">old c</div>`);
});

test('handleFetchActionErrorFields', () => {
  const elForm = createElementFromHTML<HTMLElement>(`<form>
<div class="error field"></div>
<div class="field"><input name="Foo_Bar[]"></div>
</form>`);
  handleFetchActionErrorFields(elForm, ['foo-BAR', 'other']);
  expect(normalizeTestHtml(elForm.outerHTML)).toEqual(normalizeTestHtml(`<form>
<div class="field"></div>
<div class="field error"><input name="Foo_Bar[]"></div>
</form>`));
});
