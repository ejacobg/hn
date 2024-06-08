package export

import (
	"strings"
	"testing"

	"github.com/ejacobg/hn/internal/auth"
	"github.com/ejacobg/hn/internal/item"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestFromSubmission(t *testing.T) {
	tests := []struct {
		Name string
		HTML string
		item.Story
	}{
		{
			Name: "Story",
			HTML: `<tr class='athing' id='34199828'>
                        <td align="right" valign="top" class="title"><span class="rank">1.</span></td>
                        <td valign="top" class="votelinks">
                            <center><a id='up_34199828' class='clicky nosee'
                                       href=''>
                                <div class='votearrow' title='upvote'></div>
                            </a></center>
                        </td>
                        <td class="title"><span class="titleline"><a
                                href="https://docs.technotim.live/posts/homelab-services-tour-2022/">Self-Hosting and HomeLab Walkthrough</a><span
                                class="sitebit comhead"> (<a href="from?site=technotim.live"><span class="sitestr">technotim.live</span></a>)</span></span>
                        </td>
                    </tr>`,
			Story: item.Story{
				Item:       &item.Item{ID: "34199828"},
				Title:      "Self-Hosting and HomeLab Walkthrough",
				URL:        "https://docs.technotim.live/posts/homelab-services-tour-2022/",
				Discussion: auth.BaseURL + "/item?id=34199828",
			},
		},
		{
			Name: "Ask HN",
			HTML: `<tr class='athing' id='34193766'>
                        <td align="right" valign="top" class="title"><span class="rank">2.</span></td>
                        <td valign="top" class="votelinks">
                            <center><a id='up_34193766' class='clicky nosee'
                                       href=''>
                                <div class='votearrow' title='upvote'></div>
                            </a></center>
                        </td>
                        <td class="title"><span class="titleline"><a href="item?id=34193766">Ask HN: What is the most mind expanding book(s) you have read till date?</a></span>
                        </td>
                    </tr>`,
			Story: item.Story{
				Item:       &item.Item{ID: "34193766"},
				Title:      "Ask HN: What is the most mind expanding book(s) you have read till date?",
				URL:        "item?id=34193766", // URLs are kinda broken for Ask HN's, but that's ok since the discussion link should work.
				Discussion: auth.BaseURL + "/item?id=34193766",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Log("Test", test.Name)

			nodes, err := html.ParseFragment(strings.NewReader(test.HTML), &html.Node{
				Type:     html.ElementNode,
				Data:     "tbody",
				DataAtom: atom.Tbody,
			})
			if err != nil {
				t.Fatalf("Could not parse HTML: %v", err)
			}

			node := nodes[0]
			// t.Log(node.Data)
			// for _, attr := range node.Attr {
			// 	t.Log(attr.Key, attr.Val)
			// }

			story, err := fromSubmission(node)
			if err != nil {
				t.Fatalf("Could not create item: %v", err)
			}

			if story.ID != test.ID {
				t.Logf("Expected ID %s but got %s", test.ID, story.ID)
				t.Fail()
			}

			if story.Title != test.Title {
				t.Logf("Expected Title %s but got %s", test.Title, story.Title)
				t.Fail()
			}

			if story.URL != test.URL {
				t.Logf("Expected URL %s but got %s", test.URL, story.URL)
				t.Fail()
			}

			if story.Discussion != test.Discussion {
				t.Logf("Expected Discussion %s but got %s", test.Discussion, story.Discussion)
				t.Fail()
			}
		})
	}
}

func TestFromComment(t *testing.T) {
	tests := []struct {
		Name       string
		TextLength int
		HTML       string
		item.Comment
	}{
		{
			Name:       "Top-Level Comment",
			TextLength: 2,
			HTML: `<tr class='athing' id='40549297'>    <td class='ind'></td><td valign="top" class="votelinks">
      <center><a id='up_40549297' class='clicky' href='vote?id=40549297&amp;how=up&amp;auth=2752e65950ea163f904b8e53ec476781578498d3&amp;goto=favorites%3Fid%3Dquincinia%26comments%3Dt#40549297'><div class='votearrow' title='upvote'></div></a></center>    </td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=ai_what" class="hnuser">ai_what</a> <span class="age" title="2024-06-01T21:41:02"><a href="item?id=40549297">6 days ago</a></span> <span id="unv_40549297"></span>          <span class='navs'>
             | <a href="item?id=40548807">parent</a> | <a href="context?id=40549297" rel="nofollow">context</a> | <a href="fave?id=40549297&amp;un=t&amp;auth=2752e65950ea163f904b8e53ec476781578498d3">un&#8209;favorite</a><span class="onstory"> |  on: <a href="item?id=40548807">Ask HN: Best way to save text from dynamically gen...</a></span>          </span>
                  </span></div><br><div class="comment">
                  <div class="commtext c00">You could write a userscript for this, or go a step further and write a small browser extension.<p>If you use Joplin, there&#x27;s already a web to markdown plugin for browsers.</div>
              <div class='reply'></div></div></td></tr>`,
			Comment: item.Comment{
				Item:       &item.Item{ID: "40549297"},
				Story:      "Ask HN: Best way to save text from dynamically gen...",
				Context:    auth.BaseURL + "/context?id=40549297",
				Discussion: auth.BaseURL + "/item?id=40548807",
			},
		},
		{
			Name:       "Thread Comment",
			TextLength: 1,
			HTML: `<tr class='athing' id='40255470'>    <td class='ind'></td><td valign="top" class="votelinks">
      <center><a id='up_40255470' class='clicky' href='vote?id=40255470&amp;how=up&amp;auth=8b9a3141ee3086e03ba49a7f325f926d06a39bc8&amp;goto=favorites%3Fid%3Dquincinia%26comments%3Dt#40255470'><div class='votearrow' title='upvote'></div></a></center>    </td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=rozenmd" class="hnuser">rozenmd</a> <span class="age" title="2024-05-04T06:43:49"><a href="item?id=40255470">34 days ago</a></span> <span id="unv_40255470"></span>          <span class='navs'>
             | <a href="item?id=40255457">parent</a> | <a href="context?id=40255470" rel="nofollow">context</a> | <a href="fave?id=40255470&amp;un=t&amp;auth=8b9a3141ee3086e03ba49a7f325f926d06a39bc8">un&#8209;favorite</a><span class="onstory"> |  on: <a href="item?id=40255209">A love letter to bicycle maintenance and repair</a></span>          </span>
                  </span></div><br><div class="comment">
                  <div class="commtext c00">When it&#x27;s broken, fix your bike - because when you&#x27;re broken, your bike will fix you.</div>
              <div class='reply'></div></div></td></tr>`,
			Comment: item.Comment{
				Item:       &item.Item{ID: "40255470"},
				Story:      "A love letter to bicycle maintenance and repair",
				Context:    auth.BaseURL + "/context?id=40255470",
				Discussion: auth.BaseURL + "/item?id=40255209",
			},
		},
		{
			Name:       "Multi-Paragraph Comment",
			TextLength: 4,
			HTML: `<tr class='athing' id='32763668'>    <td class='ind'></td><td valign="top" class="votelinks">
      <center><a id='up_32763668' class='clicky' href='vote?id=32763668&amp;how=up&amp;auth=f07e8ae8b12b97f6f3d25eab234cf6c1f83d9d66&amp;goto=favorites%3Fid%3Dquincinia%26comments%3Dt#32763668'><div class='votearrow' title='upvote'></div></a></center>    </td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=nicbou" class="hnuser">nicbou</a> <span class="age" title="2022-09-08T10:34:55"><a href="item?id=32763668">on Sept 8, 2022</a></span> <span id="unv_32763668"></span>          <span class='navs'>
             | <a href="item?id=32746922">parent</a> | <a href="context?id=32763668" rel="nofollow">context</a> | <a href="fave?id=32763668&amp;un=t&amp;auth=f07e8ae8b12b97f6f3d25eab234cf6c1f83d9d66">un&#8209;favorite</a><span class="onstory"> |  on: <a href="item?id=32746922">Excuse me but why are you eating so many frogs</a></span>          </span>
                  </span></div><br><div class="comment">
                  <div class="commtext c00">The topic is nice, but above all I thoroughly enjoy the writing.<p>I feel like I get a lot more done precisely because I don&#x27;t like to eat frogs. I sleep until I&#x27;m rested. I don&#x27;t touch my computer until I&#x27;ve had tea on the balcony. I work on what feels right, when it feels right, for as long as it feels right. If the weather is nice, I&#x27;ll hop on my bicycle and forget about work.<p>But when something sparks my interest, I have stores of energy to throw at it. My appetite for work is unrestrained by the frogs I&#x27;ve had for breakfast.<p>I embraced the fact that I am not a machine, and that my output is neither constant nor predictable. I&#x27;d rather respect the tides of my energy than fight against them.</div>
              <div class='reply'></div></div></td></tr>`,
			Comment: item.Comment{
				Item:       &item.Item{ID: "32763668"},
				Story:      "Excuse me but why are you eating so many frogs",
				Context:    auth.BaseURL + "/context?id=32763668",
				Discussion: auth.BaseURL + "/item?id=32746922",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Log("Test", test.Name)

			nodes, err := html.ParseFragment(strings.NewReader(test.HTML), &html.Node{
				Type:     html.ElementNode,
				Data:     "tbody",
				DataAtom: atom.Tbody,
			})
			if err != nil {
				t.Fatalf("Could not parse HTML: %v", err)
			}

			node := nodes[0]
			// t.Log(node.Data)
			// for _, attr := range node.Attr {
			// 	t.Log(attr.Key, attr.Val)
			// }

			comment, err := fromComment(node)
			if err != nil {
				t.Fatalf("Could not create item: %v", err)
			}

			if comment.ID != test.ID {
				t.Logf("Expected ID %s but got %s", test.ID, comment.ID)
				t.Fail()
			}

			// t.Logf("%#v", comment.Text)
			if len(comment.Text) != test.TextLength {
				t.Logf("Expected %d paragraphs but got %d", test.TextLength, len(comment.Text))
				t.Fail()
			}

			if comment.Story != test.Story {
				t.Logf("Expected Story %s but got %s", test.Story, comment.Story)
				t.Fail()
			}

			if comment.Context != test.Context {
				t.Logf("Expected Context %s but got %s", test.Context, comment.Context)
				t.Fail()
			}

			if comment.Discussion != test.Discussion {
				t.Logf("Expected Discussion %s but got %s", test.Discussion, comment.Discussion)
				t.Fail()
			}
		})
	}
}
