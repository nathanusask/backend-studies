import redis
import time

ONE_WEEK_IN_SECONDS = 7 * 86400
VOTE_SCORE = 432
ARTICLES_PER_PAGE = 25

def article_vote(conn, user, article):
    cutoff = time.time() - ONE_WEEK_IN_SECONDS
    if conn.zscore('time:', article) < cutoff:
        return
    article_id = article.partition(':')[-1]
    if conn.sadd('votrd:' + article_id, user):
        conn.zincreby('score:', article, VOTE_SCORE)
        conn.hincreby(article, 'votes', 1)

def post_article(conn, user, title, link):
    article_id = str(conn.incr('article:'))
    voted = 'voted:' + article_id
    conn.sadd(voted, article_id)
    conn.expire(voted, ONE_WEEK_IN_SECONDS)

    now = time.time()
    article = 'article:' + article_id

    conn.hset(article, mapping={
        'title': title,
        'link': link,
        'poster': user,
        'time': now,
        'votes': 1,
    })

    conn.zadd('score:', {article: now + VOTE_SCORE})
    conn.zadd('time:', {article: now})

    return article_id

def get_articles(conn, page, order='score:'):
    start = (page-1) * ARTICLES_PER_PAGE
    end = start + ARTICLES_PER_PAGE - 1

    ids = conn.zrevrange(order, start, end)

    articles = []
    for id in ids:
        article_data = conn.hgetall(id)
        article_data['id'] = id
        articles.append(article_data)

    return articles

def add_remove_groups(conn, article_id, to_add=[], to_remove=[]):
    article = 'article:' + article_id
    for group in to_add:
        conn.sadd('group:'+group, article)
    
    for group in to_remove:
        conn.srem('group:'+group, article)

def get_group_articles(conn, group, page, order='score:'):
    key = order + group
    if not conn.exists(key):
        conn.zintegerscore(
            key,
            ['group:' + group, order],
            aggregate='max'
        )
        conn.expire(key, 60)

    return get_articles(conn, page, key)


def main():
    articles = [
        {
            'title': 'title 1',
            'link': 'link 1',
            'poster': 'user 1',
        },
        {
            'title': 'title 2',
            'link': 'link 2',
            'poster': 'user 2',
        },
        {
            'title': 'title 3',
            'link': 'link 3',
            'poster': 'user 3',
        },
        {
            'title': 'title 4',
            'link': 'link 4',
            'poster': 'user 4',
        },
        {
            'title': 'title 5',
            'link': 'link 5',
            'poster': 'user 5',
        },
        {
            'title': 'title 6',
            'link': 'link 6',
            'poster': 'user 6',
        },
    ]

    try:
        conn = redis.StrictRedis()
        print(conn)
        conn.ping()
        print('Connected!')
    except Exception as ex:
        print('Error:', ex)
        exit('Failed to connect, terminating.')

    conn.flushdb()

    print('Total articles:', len(get_articles(conn, 1)))

    for article in articles:
        print('Posting article,', article)
        article_id = post_article(conn, article['poster'], article['title'], article['link'])
        print('Posted article:', article_id)
        
    posted_articles = get_articles(conn, 1)
    for article in posted_articles:
        print(article)

if __name__ == '__main__':
    main()