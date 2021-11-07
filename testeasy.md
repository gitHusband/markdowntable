<table style="width:100%">
<thead>
	<th>参数</th>
	<th>子参</th>
	<th>子参</th>
	<th>释义</th>
</thead>
<tbody>
	<tr>
		<td>baseElement</td>
		<td colspan="3"><strong style="font-size: 15px">Tips 基准元素 </strong><br/> <em style="color: #888888">Tips 定位的参照元素</em><br/> <b>默认：<ins>空值</ins></b><br/> <b>可能的值：</b> <ul><li>空值 - 代表Body元素</li><li>DOM 或 jQuery 元素 </li><li>DOM 选择器</li></ul></td>
	</tr>
	<tr>
		<td rowspan="2">offset</td>
		<td>left</td>
		<td colspan="2"><strong style="font-size: 15px">T偏移父元素Left 多少px </strong><br/> <em style="color: #888888">任意实数</em></td>
	</tr>
	<tr>
		<td>top</td>
		<td colspan="2"><strong style="font-size: 15px">偏移父元素Top 多少px</strong><br/> <em style="color: #888888">任意实数</em></td>
	</tr>
	<tr>
		<td rowspan="4">symbolOptions</td>
		<td rowspan="2">offset</td>
		<td>left</td>
		<td><strong style="font-size: 15px">偏移Tips元素Top 多少px(暂不支持) </strong><br/> <em style="color: #888888">任意实数</em></td>
	</tr>
	<tr>
		<td>top</td>
		<td><strong style="font-size: 15px">偏移Tips元素Top 多少px(暂不支持) </strong><br/> <em style="color: #888888">任意实数</em></td>
	</tr>
	<tr>
		<td>position</td>
		<td colspan="2"><strong style="font-size: 15px">symbol与Tips的相对定位</strong><br/> <b>默认：<ins>top-center</ins></b><br/> <b>可能的值：</b> <ul><li>top-left</li><li>top-center</li><li>top-right </li></ul></td>
	</tr>
	<tr>
		<td>type</td>
		<td colspan="2"><strong style="font-size: 15px">如果不设置，symbol则按option.type 色调，否则按该参数的色调</strong><br/> <b>默认：<ins>info</ins></b><br/> <b>可能的值：</b> <ul><li>normal</li><li>success</li><li>error</li></ul></td>
	</tr>
</tbody>
</table>
