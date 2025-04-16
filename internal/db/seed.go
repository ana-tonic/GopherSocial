package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/ana-tonic/gopher-social/internal/store"
)

var usernames = []string{
	"Alice", "Bob", "Carla", "David", "Emma",
	"Frank", "Grace", "Harry", "Isla", "Jack",
	"Kelly", "Liam", "Mia", "Nathan", "Olivia",
	"Peter", "Quinn", "Rachel", "Sam", "Tina",
	"Ursula", "Victor", "Wendy", "Xander", "Yasmin",
	"Zach", "Amy", "Bradley", "Chloe", "Derek",
	"Ella", "Felix", "Gina", "Hank", "Ivy",
	"Jason", "Karen", "Leo", "Melanie", "Noah",
	"Opal", "Paul", "Ruby", "Ryan", "Sophia",
	"Tom", "Una", "Vicky", "Will", "Zoe",
}

var titles = []string{
	"Why Simplicity Wins",
	"Mastering the Basics",
	"Go Beyond the Code",
	"Small Steps, Big Impact",
	"The Art of Focus",
	"Build Without Burnout",
	"Rethinking Productivity",
	"Lessons from Failure",
	"Design for Humans",
	"Code with Clarity",
	"Fast Doesn’t Mean Rushed",
	"Debugging Your Mind",
	"Start Before You’re Ready",
	"Make It Make Sense",
	"Habits That Stick",
	"Write Fewer Features",
	"Think Like a Beginner",
	"Clean Code, Clean Mind",
	"From Idea to Action",
	"Don’t Wait to Launch",
}

var contents = []string{
	"Sometimes the simplest solution is the most powerful. Focus on clarity, not complexity.",
	"Learning the basics isn't glamorous, but it's the foundation of everything you build.",
	"Great code starts with great thinking. Write it like someone else has to read it tomorrow.",
	"Every small improvement compounds over time. Stay consistent.",
	"Distraction is the enemy of progress. Cut the noise and focus on what matters.",
	"Burnout comes from imbalance. Work with intent, rest with purpose.",
	"Productivity isn’t about doing more. It’s about doing what matters most.",
	"Failure is feedback in disguise. Learn, adapt, and move forward.",
	"Design is empathy made visible. Build with the user in mind.",
	"Readable code is maintainable code. Favor clarity over cleverness.",
	"Speed matters, but quality matters more. Ship fast, but ship smart.",
	"Debugging starts with understanding. Get curious about the bug.",
	"Perfection is the enemy of momentum. Start scrappy and refine as you go.",
	"If it’s not clear, it’s not done. Communicate your ideas fully.",
	"Habits shape outcomes. Build daily rituals that support your goals.",
	"Not every feature is worth building. Sometimes less really is more.",
	"Approach problems like it’s your first time. Stay curious.",
	"Clean code reflects clean thinking. Refactor with care.",
	"Turn your ideas into experiments. Progress beats perfection.",
	"Done is better than perfect. Publish, then improve.",
}

var tags = []string{
	"productivity",
	"mindset",
	"coding",
	"golang",
	"design",
	"startups",
	"leadership",
	"creativity",
	"focus",
	"habits",
	"clean-code",
	"learning",
	"motivation",
	"software-development",
	"ui-ux",
	"debugging",
	"remote-work",
	"growth",
	"minimalism",
	"writing",
}

var comments = []string{
	"Great post! Thanks for sharing.",
	"I completely agree with your thoughts.",
	"Thanks for the tips, very helpful.",
	"Interesting perspective, I hadn't considered that.",
	"Thanks for sharing your experience.",
	"Well written, I enjoyed reading this.",
	"This is very insightful, thanks for posting.",
	"Great advice, I'll definitely try that.",
	"I love this, very inspirational.",
	"Thanks for the information, very useful.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			tx.Rollback()
			log.Println("Error creating user", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Println("Error creating post", err)
			return
		}
	}

	comments := generateComments(500, posts, users)

	for _, comment := range comments {
		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Println("Error creating comment", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))]},
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {

		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
