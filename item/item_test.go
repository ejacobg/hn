package item

import (
	"github.com/ejacobg/hn/auth"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"testing"
)

func TestFromSubmission(t *testing.T) {
	tests := []struct {
		Name string
		HTML string
		Story
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
			Story: Story{
				ID:         "34199828",
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
			Story: Story{
				ID:         "34193766",
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

			story, err := FromSubmission(node)
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
		Comment
	}{
		{
			Name:       "Top-Level Comment",
			TextLength: 5,
			HTML: `<tr class='athing' id='34210369'>
                        <td class='ind'></td>
                        <td valign="top" class="votelinks">
                            <center><a id='up_34210369' class='clicky nosee'
                                       href=''>
                                <div class='votearrow' title='upvote'></div>
                            </a></center>
                        </td>
                        <td class="default">
                            <div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=thyrox" class="hnuser">thyrox</a> <span class="age" title="2023-01-01T20:14:51"><a
                                    href="item?id=34210369">4 days ago</a></span> <span id="unv_34210369"></span>          <span
                                    class='navs'>
             | <a href="item?id=34206219">parent</a> | <a href="context?id=34210369" rel="nofollow">context</a> | <a
                                    href="upvoted?id=quincinia&amp;comments=t&amp;p=2" aria-hidden="true">next</a><span
                                    class="onstory"> |  on: <a href="item?id=34206219">Ask HN: Concepts that clicked only years after you...</a></span>          </span>
                  </span></div>
                            <br>
                            <div class="comment">
                                <span class="commtext c00">The power of follow ups (especially in sales)<p>One thing which held me back for a very long time was not following up with people who didn&#x27;t show much interest initially.<p>I wasted so many good leads thinking it is impolite to follow up with people after contacting them once. My whole life changed once I understood the power of follow ups and understanding that most people are so busy that it takes at least 6 reminders before most people will take any substantial action.<p>The reverse is also true. People say a lot of things and most of the times you never cross the bridge or reach it. Nowadays, I rarely argue about anything and don&#x27;t act on stuff until a person reminds me once or twice. This small filter can be like a miracle for saving your time and energy.</span>
                                <div class='reply'></div>
                            </div>
                        </td>
                    </tr>`,
			Comment: Comment{
				ID:         "34210369",
				Story:      "Ask HN: Concepts that clicked only years after you...",
				Context:    auth.BaseURL + "/context?id=34210369",
				Discussion: auth.BaseURL + "/item?id=34206219",
			},
		},
		{
			Name:       "Thread Comment",
			TextLength: 3,
			HTML: `<tr class='athing' id='34188913'>
                        <td class='ind'></td>
                        <td valign="top" class="votelinks">
                            <center><a id='up_34188913' class='clicky nosee'
                                       href=''>
                                <div class='votearrow' title='upvote'></div>
                            </a></center>
                        </td>
                        <td class="default">
                            <div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=nnadams" class="hnuser">nnadams</a> <span class="age" title="2022-12-30T19:27:03"><a
                                    href="item?id=34188913">6 days ago</a></span> <span id="unv_34188913"></span>          <span
                                    class='navs'>
             | <a href="item?id=34188687">parent</a> | <a href="context?id=34188913" rel="nofollow">context</a> | <a
                                    href="upvoted?id=quincinia&amp;comments=t&amp;p=2" aria-hidden="true">next</a><span
                                    class="onstory"> |  on: <a
                                    href="item?id=34186283">Why I'm still using Python</a></span>          </span>
                  </span></div>
                            <br>
                            <div class="comment">
                                <span class="commtext c00">I have embraced this workflow completely. I used to be more concerned about which language, but now I find I much more useful to just start immediately in Python. I spend most of the time working out the kinks and edge cases, instead of memory management or other logistics. Maybe ~75% of the time, that&#x27;s it, no need for further improvement.<p>Recently I chose to rewrite several thousands of lines of Python in Go, because we needed more speed and improved concurrency. Already having a working program and tests in Python was great. After figuring out a few Go-isms, it was a quick couple of days to port it all for big improvements.</span>
                                <div class='reply'></div>
                            </div>
                        </td>
                    </tr>`,
			Comment: Comment{
				ID:         "34188913",
				Story:      "Why I'm still using Python",
				Context:    auth.BaseURL + "/context?id=34188913",
				Discussion: auth.BaseURL + "/item?id=34186283",
			},
		},
		{
			Name:       "1-Paragraph Comment",
			TextLength: 1,
			HTML: `<tr class='athing' id='34148489'>
                        <td class='ind'></td>
                        <td valign="top" class="votelinks">
                            <center><a id='up_34148489' class='clicky nosee'
                                       href=''>
                                <div class='votearrow' title='upvote'></div>
                            </a></center>
                        </td>
                        <td class="default">
                            <div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=synu" class="hnuser">synu</a> <span class="age" title="2022-12-27T13:40:51"><a
                                    href="item?id=34148489">9 days ago</a></span> <span id="unv_34148489"></span>          <span
                                    class='navs'>
             | <a href="item?id=34145680">parent</a> | <a href="context?id=34148489" rel="nofollow">context</a> | <a
                                    href="upvoted?id=quincinia&amp;comments=t&amp;p=2" aria-hidden="true">next</a><span
                                    class="onstory"> |  on: <a
                                    href="item?id=34145680">I have reached Vim nirvana</a></span>          </span>
                  </span></div>
                            <br>
                            <div class="comment">
                                <span class="commtext c00">Vim enlightenment for me was learning to use it with all the default settings. Everywhere I go, it’s already set up just the way I like it.</span>
                                <div class='reply'></div>
                            </div>
                        </td>
                    </tr>`,
			Comment: Comment{
				ID:         "34148489",
				Story:      "I have reached Vim nirvana",
				Context:    auth.BaseURL + "/context?id=34148489",
				Discussion: auth.BaseURL + "/item?id=34145680",
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

			comment, err := FromComment(node)
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
