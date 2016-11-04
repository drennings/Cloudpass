

import java.io.IOException;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Scanner;

import javax.servlet.ServletException;
import javax.servlet.annotation.MultipartConfig;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.Part;

/**
 * Servlet implementation class UploadServlet
 */
@WebServlet("/UploadServlet")
@MultipartConfig(fileSizeThreshold=1024*1024*2, // 2MB
                 maxFileSize=1024*1024*10,      // 10MB
                 maxRequestSize=1024*1024*50)   // 50MB
public class UploadServlet extends HttpServlet {
	private static final long serialVersionUID = 1L;
       

        /**
         * Name of the directory where uploaded files will be saved, relative to
         * the web application directory.
         */
        private static final String SAVE_DIR = "uploadFiles";
        
     // JDBC driver name and database URL
    	   static final String JDBC_DRIVER = "com.mysql.jdbc.Driver";  
    	   static final String DB_URL = "jdbc:mysql://cloudpassdb.cj7d5cmvkhd4.us-west-2.rds.amazonaws.com:3306/cloudpass";

    	   //  Database credentials
    	   static final String USER = "cpmaster";
    	   static final String PASS = "cloudpass939";
         
        /**
         * handles file upload
         */
        protected void doPost(HttpServletRequest request,
                HttpServletResponse response) throws ServletException, IOException {
        	System.out.println("start doPost");
        	
        	//Extract parameters
        	String hash = extractHash(request.getPart("file"));
        	String user = extractHash(request.getPart("user"));
        	String email = extractHash(request.getPart("email"));
        	String capacity = extractHash(request.getPart("capacity"));
        	String htype = extractHash(request.getPart("htype"));    
            System.out.println("Hash: "+hash);
            System.out.println("User: "+user);
            System.out.println("Email: "+email);
            System.out.println("Capacity: "+capacity);
            System.out.println("Hash type: "+htype);
        	
        	//connect    	    
        	   Connection conn = null;
        	   Statement stmt = null;
        	   try{
        	      //STEP 2: Register JDBC driver
        	      Class.forName("com.mysql.jdbc.Driver");

        	      //STEP 3: Open a connection
        	      System.out.println("Connecting to database...");
        	      conn = DriverManager.getConnection(DB_URL,USER,PASS);
        	      
        	    //STEP 4: Execute a query
        	      System.out.println("Creating table in given database...");
        	      stmt = conn.createStatement();
        	      
        	      String sql = "CREATE TABLE if not exists JOBS " +
        	                   "(id INTEGER not NULL, " +
        	                   " hash VARCHAR(255), " + 
        	                   " user VARCHAR(255), " + 
        	                   " email VARCHAR(255), " + 
        	                   " capacity INTEGER, " +
        	                   " htype VARCHAR(255), " + 
        	                   " PRIMARY KEY ( id ))"; 

        	      stmt.executeUpdate(sql);
        	      System.out.println("Created table in given database...");
        	      
        	      sql = "INSERT INTO JOBS " +
        	    		  "VALUES ("+"001, "+"'"+hash+"', '"+user+"', '"+email+"', "+capacity+", '"+htype+"')";
        	      stmt.executeUpdate(sql);
        	   } catch (Exception e) {
       			// TODO Auto-generated catch block
       			e.printStackTrace();
       		}finally{
       	      //finally block used to close resources
       	      try{
       	         if(stmt!=null)
       	            conn.close();
       	      }catch(SQLException se){
       	      }// do nothing
       	      try{
       	         if(conn!=null)
       	            conn.close();
       	      }catch(SQLException se){
       	         se.printStackTrace();
       	      }//end finally try
       		}
        	             	
        	     
     
            request.setAttribute("message", "Upload has been done successfully!");
            getServletContext().getRequestDispatcher("/message.jsp").forward(
                    request, response);
        }
        
        /**
         * Extracts file content from HTTP header content-disposition
         */
        private String extractHash(Part part) {
        	String hash = "No hash found";
    		try {
    			Scanner scan = new Scanner(part.getInputStream());
    			hash = scan.useDelimiter("\\Z").next();
    		} catch (IOException e) {
    			// TODO Auto-generated catch block
    			e.printStackTrace();
    		}
            return hash;
        }
        

}
