// MongoDB initialization script
// This script runs automatically when MongoDB container starts

db = db.getSiblingDB('message_dispatcher');

// Create messages collection
db.createCollection('messages');

// Insert sample test messages
db.messages.insertMany([
    {
        phone_number: "+905551111111",
        content: "Test mesajı 1 - Insider Project",
        status: "pending",
        created_at: new Date()
    },
    {
        phone_number: "+905552222222",
        content: "Test mesajı 2 - Welcome to Insider",
        status: "pending",
        created_at: new Date()
    },
    {
        phone_number: "+905553333333",
        content: "Test mesajı 3 - Siparişiniz hazır",
        status: "pending",
        created_at: new Date()
    },
    {
        phone_number: "+905554444444",
        content: "Test mesajı 4 - Alışveriş için teşekkürler",
        status: "pending",
        created_at: new Date()
    },
    {
        phone_number: "+905555555555",
        content: "Test mesajı 5 - Sizin için özel indirim",
        status: "pending",
        created_at: new Date()
    }
]);

print("MongoDB initialized successfully!");
print("Collection 'messages' created");
print("5 sample messages inserted");

